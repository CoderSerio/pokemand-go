package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const defaultPreviewLines = 5

type ScriptInfo struct {
	Name         string    `json:"name"`
	RelativePath string    `json:"relativePath"`
	AbsolutePath string    `json:"absolutePath"`
	SizeBytes    int64     `json:"sizeBytes"`
	Mode         string    `json:"mode"`
	Executable   bool      `json:"executable"`
	Shebang      string    `json:"shebang,omitempty"`
	ModifiedAt   time.Time `json:"modifiedAt"`
	Preview      []string  `json:"preview,omitempty"`
}

func getScriptsDir() string {
	return filepath.Join(GetDataDir(), "scripts")
}

func ensureScriptsDir() (string, error) {
	scriptsDir := getScriptsDir()
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return "", err
	}
	return scriptsDir, nil
}

func listManagedScripts(previewLines int) ([]ScriptInfo, error) {
	scriptsDir, err := ensureScriptsDir()
	if err != nil {
		return nil, fmt.Errorf("创建脚本目录失败: %w", err)
	}

	var scripts []ScriptInfo
	err = filepath.WalkDir(scriptsDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		info, err := buildScriptInfo(scriptsDir, path, previewLines)
		if err != nil {
			return err
		}
		scripts = append(scripts, info)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("读取脚本目录失败: %w", err)
	}

	sort.Slice(scripts, func(i, j int) bool {
		return scripts[i].RelativePath < scripts[j].RelativePath
	})

	return scripts, nil
}

func findManagedScript(path string, previewLines int) (ScriptInfo, error) {
	scriptsDir, err := ensureScriptsDir()
	if err != nil {
		return ScriptInfo{}, fmt.Errorf("创建脚本目录失败: %w", err)
	}

	candidates := []string{
		path,
		filepath.Join(scriptsDir, path),
	}

	for _, candidate := range candidates {
		fileInfo, err := os.Stat(candidate)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return ScriptInfo{}, err
		}
		if fileInfo.IsDir() {
			return ScriptInfo{}, fmt.Errorf("目标是目录，不是文件: %s", candidate)
		}

		return buildScriptInfo(scriptsDir, candidate, previewLines)
	}

	return ScriptInfo{}, os.ErrNotExist
}

func filterScriptsBySearch(scripts []ScriptInfo, search string) []ScriptInfo {
	if search == "" {
		return scripts
	}

	var filtered []ScriptInfo
	needle := strings.ToLower(search)
	for _, script := range scripts {
		if strings.Contains(strings.ToLower(script.RelativePath), needle) ||
			strings.Contains(strings.ToLower(script.Name), needle) {
			filtered = append(filtered, script)
		}
	}

	return filtered
}

func printJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func buildScriptInfo(baseDir string, path string, previewLines int) (ScriptInfo, error) {
	resolvedPath, err := filepath.Abs(path)
	if err != nil {
		return ScriptInfo{}, err
	}

	fileInfo, err := os.Stat(resolvedPath)
	if err != nil {
		return ScriptInfo{}, err
	}

	relativePath, err := filepath.Rel(baseDir, resolvedPath)
	if err != nil {
		return ScriptInfo{}, err
	}

	var (
		preview []string
		shebang string
	)
	if previewLines > 0 {
		preview, shebang, err = readPreview(resolvedPath, previewLines)
		if err != nil {
			return ScriptInfo{}, err
		}
	}

	return ScriptInfo{
		Name:         filepath.Base(resolvedPath),
		RelativePath: filepath.ToSlash(relativePath),
		AbsolutePath: resolvedPath,
		SizeBytes:    fileInfo.Size(),
		Mode:         fileInfo.Mode().String(),
		Executable:   fileInfo.Mode().Perm()&0111 != 0,
		Shebang:      shebang,
		ModifiedAt:   fileInfo.ModTime(),
		Preview:      preview,
	}, nil
}

func readPreview(path string, maxLines int) ([]string, string, error) {
	if maxLines <= 0 {
		maxLines = defaultPreviewLines
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var preview []string
	var shebang string

	for len(preview) < maxLines {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, "", err
		}

		line = strings.TrimRight(line, "\r\n")
		if shebang == "" && strings.HasPrefix(line, "#!") {
			shebang = line
		}
		if line != "" || err != io.EOF {
			preview = append(preview, line)
		}

		if err == io.EOF {
			break
		}
	}

	return preview, shebang, nil
}
