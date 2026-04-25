package cmd

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

//go:embed webui/index.html
var webUIFS embed.FS

type wsMessage struct {
	Version string          `json:"version,omitempty"`
	Type    string          `json:"type"`
	ID      string          `json:"id,omitempty"`
	Action  string          `json:"action,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
	OK      bool            `json:"ok,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Error   *wsError        `json:"error,omitempty"`
}

type wsError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type uiServer struct {
	token    string
	upgrader websocket.Upgrader
}

type bootstrapData struct {
	App struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"app"`
	Local struct {
		ScriptsDir string `json:"scriptsDir"`
	} `json:"local"`
	Market struct {
		Provider   string `json:"provider"`
		Configured bool   `json:"configured"`
		DocsURL    string `json:"docsUrl"`
		BaseURL    string `json:"baseUrl"`
	} `json:"market"`
}

var (
	uiHost   string
	uiPort   int
	uiNoOpen bool
)

var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the local Web UI",
	Long:  "Start the local Web UI for managing skill scripts over a WebSocket-backed session",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runUIServer(); err != nil {
			fmt.Printf("Failed to start UI server: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().StringVar(&uiHost, "host", "127.0.0.1", "Host address to bind the UI server to")
	uiCmd.Flags().IntVar(&uiPort, "port", 3719, "Port for the UI server, or 0 to auto-assign")
	uiCmd.Flags().BoolVar(&uiNoOpen, "no-open", false, "Do not open the browser automatically")
}

func runUIServer() error {
	token, err := randomToken()
	if err != nil {
		return err
	}

	server := &uiServer{
		token: token,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleIndex)
	mux.HandleFunc("/ws", server.handleWS)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", uiHost, uiPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	httpServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	url := fmt.Sprintf("http://%s", listener.Addr().String())
	fmt.Printf("pkmg ui is running at %s\n", url)

	if !uiNoOpen {
		if err := openSystemPath(url); err != nil {
			fmt.Printf("Failed to open browser automatically: %v\n", err)
		}
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpServer.Serve(listener)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case sig := <-sigCh:
		fmt.Printf("Received signal %s, shutting down UI...\n", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return httpServer.Shutdown(ctx)
}

func (s *uiServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	content, err := webUIFS.ReadFile("webui/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := strings.ReplaceAll(string(content), "__PKMG_WS_TOKEN__", s.token)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, html)
}

func (s *uiServer) handleWS(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("token") != s.token {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		var msg wsMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return
		}
		if msg.Type != "request" {
			_ = conn.WriteJSON(wsMessage{
				Version: "1.0",
				Type:    "response",
				ID:      msg.ID,
				OK:      false,
				Error:   &wsError{Code: "invalid_type", Message: "message type must be request"},
			})
			continue
		}

		data, err := handleUIAction(msg.Action, msg.Payload)
		resp := wsMessage{
			Version: "1.0",
			Type:    "response",
			ID:      msg.ID,
			OK:      err == nil,
			Data:    data,
		}
		if err != nil {
			resp.Error = &wsError{Code: "action_failed", Message: err.Error()}
		}
		if writeErr := conn.WriteJSON(resp); writeErr != nil {
			return
		}
	}
}

func handleUIAction(action string, payload json.RawMessage) (interface{}, error) {
	switch action {
	case "app.bootstrap":
		return getBootstrapData(), nil
	case "skills.list":
		var req struct {
			Query string `json:"query"`
		}
		_ = json.Unmarshal(payload, &req)
		if req.Query != "" {
			return searchLocalSkills(req.Query)
		}
		return listLocalSkills()
	case "skills.get":
		var req struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return getLocalSkillDetail(req.ID)
	case "skills.create":
		var req struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return createLocalSkillWithRequiredExtension(req.Path, req.Content)
	case "skills.save":
		var req struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return saveLocalSkill(req.ID, req.Content)
	case "skills.copy":
		var req struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return copyLocalSkill(req.ID)
	case "skills.delete":
		var req struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		if err := deleteLocalSkill(req.ID); err != nil {
			return nil, err
		}
		return map[string]string{"status": "deleted"}, nil
	case "skills.restoreVersion":
		var req struct {
			ID      string `json:"id"`
			Version int    `json:"version"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return restoreLocalSkillVersion(req.ID, req.Version)
	case "skills.openDir":
		var req struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		if err := openLocalSkillDir(req.ID); err != nil {
			return nil, err
		}
		return map[string]string{"status": "ok"}, nil
	case "market.catalog":
		return fetchSkillHubCatalog(12)
	case "market.search":
		var req struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal(payload, &req); err != nil {
			return nil, err
		}
		return fetchSkillHubSearch(req.Query, 12)
	default:
		return nil, fmt.Errorf("unsupported action: %s", action)
	}
}

func getBootstrapData() bootstrapData {
	var data bootstrapData
	data.App.Name = "pkmg"
	data.App.Version = Version
	data.Local.ScriptsDir = getScriptsDir()
	data.Market.Provider = "SkillHub"
	data.Market.Configured = strings.TrimSpace(os.Getenv("SKILLHUB_API_KEY")) != ""
	data.Market.DocsURL = "https://www.skillhub.club/docs/api"
	data.Market.BaseURL = "https://www.skillhub.club/api/v1"
	return data
}

func randomToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
