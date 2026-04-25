package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SkillVersion struct {
	Number   int       `json:"number"`
	Label    string    `json:"label"`
	SavedAt  time.Time `json:"savedAt"`
	Note     string    `json:"note,omitempty"`
	FileName string    `json:"fileName"`
}

type skillMetadata struct {
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	NextVersion int            `json:"nextVersion"`
	Versions    []SkillVersion `json:"versions"`
}

type LocalSkill struct {
	ID           string         `json:"id"`
	Title        string         `json:"title"`
	Summary      string         `json:"summary"`
	RelativePath string         `json:"relativePath"`
	AbsolutePath string         `json:"absolutePath"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	VersionCount int            `json:"versionCount"`
	Versions     []SkillVersion `json:"versions"`
}

type LocalSkillDetail struct {
	LocalSkill
	Content  string `json:"content"`
	Shebang  string `json:"shebang,omitempty"`
	FileMode string `json:"fileMode"`
}

func getSkillStateRoot() string {
	return filepath.Join(GetDataDir(), ".pkmg")
}

func getSkillMetadataPath(relativePath string) string {
	return filepath.Join(getSkillStateRoot(), "skills", filepath.FromSlash(relativePath)+".json")
}

func getSkillVersionsDir(relativePath string) string {
	return filepath.Join(getSkillStateRoot(), "versions", filepath.FromSlash(relativePath))
}

func ensureSkillStateDirs(relativePath string) error {
	paths := []string{
		filepath.Dir(getSkillMetadataPath(relativePath)),
		getSkillVersionsDir(relativePath),
	}
	for _, path := range paths {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func listLocalSkills() ([]LocalSkill, error) {
	scripts, err := listManagedScripts(0)
	if err != nil {
		return nil, err
	}

	skills := make([]LocalSkill, 0, len(scripts))
	for _, script := range scripts {
		skill, err := buildLocalSkill(script)
		if err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	return skills, nil
}

func searchLocalSkills(query string) ([]LocalSkill, error) {
	skills, err := listLocalSkills()
	if err != nil {
		return nil, err
	}
	if query == "" {
		return skills, nil
	}

	var results []LocalSkill
	needle := strings.ToLower(query)
	for _, skill := range skills {
		if strings.Contains(strings.ToLower(skill.Title), needle) ||
			strings.Contains(strings.ToLower(skill.Summary), needle) ||
			strings.Contains(strings.ToLower(skill.RelativePath), needle) {
			results = append(results, skill)
		}
	}

	return results, nil
}

func getLocalSkillDetail(relativePath string) (LocalSkillDetail, error) {
	script, err := findManagedScript(relativePath, defaultPreviewLines)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	skill, _, err := buildLocalSkillWithMetadata(script)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	content, err := os.ReadFile(script.AbsolutePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	return LocalSkillDetail{
		LocalSkill: skill,
		Content:    string(content),
		Shebang:    script.Shebang,
		FileMode:   script.Mode,
	}, nil
}

func saveLocalSkill(relativePath string, content string) (LocalSkillDetail, error) {
	script, err := findManagedScript(relativePath, defaultPreviewLines)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	currentContent, err := os.ReadFile(script.AbsolutePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}
	if string(currentContent) == content {
		return getLocalSkillDetail(relativePath)
	}

	metadata, err := ensureSkillMetadata(script)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	fileInfo, err := os.Stat(script.AbsolutePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	now := time.Now()
	if err := os.WriteFile(script.AbsolutePath, []byte(content), fileInfo.Mode().Perm()); err != nil {
		return LocalSkillDetail{}, err
	}

	metadata, err = appendSkillVersion(relativePath, metadata, []byte(content), now, "Saved from Web UI")
	if err != nil {
		return LocalSkillDetail{}, err
	}

	if err := writeSkillMetadata(relativePath, metadata); err != nil {
		return LocalSkillDetail{}, err
	}

	return getLocalSkillDetail(relativePath)
}

func createLocalSkill(relativePath string, content string) (LocalSkillDetail, error) {
	return createLocalSkillWithPathNormalizer(relativePath, content, normalizeSkillRelativePath)
}

func createLocalSkillWithRequiredExtension(relativePath string, content string) (LocalSkillDetail, error) {
	return createLocalSkillWithPathNormalizer(relativePath, content, normalizeSkillRelativePathWithRequiredExtension)
}

func createLocalSkillWithPathNormalizer(relativePath string, content string, normalize func(string) (string, error)) (LocalSkillDetail, error) {
	normalizedPath, err := normalize(relativePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	targetAbsolutePath := filepath.Join(getScriptsDir(), filepath.FromSlash(normalizedPath))
	if _, err := os.Stat(targetAbsolutePath); err == nil {
		return LocalSkillDetail{}, fmt.Errorf("managed skill script already exists: %s", normalizedPath)
	} else if !os.IsNotExist(err) {
		return LocalSkillDetail{}, err
	}

	if err := os.MkdirAll(filepath.Dir(targetAbsolutePath), 0755); err != nil {
		return LocalSkillDetail{}, err
	}

	if strings.TrimSpace(content) == "" {
		content = defaultSkillTemplate(normalizedPath)
	}

	if err := os.WriteFile(targetAbsolutePath, []byte(content), 0755); err != nil {
		return LocalSkillDetail{}, err
	}

	return getLocalSkillDetail(normalizedPath)
}

func copyLocalSkill(relativePath string) (LocalSkillDetail, error) {
	script, err := findManagedScript(relativePath, defaultPreviewLines)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	content, err := os.ReadFile(script.AbsolutePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	targetRelativePath, err := nextCopyRelativePath(relativePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	targetAbsolutePath := filepath.Join(getScriptsDir(), filepath.FromSlash(targetRelativePath))
	if err := os.MkdirAll(filepath.Dir(targetAbsolutePath), 0755); err != nil {
		return LocalSkillDetail{}, err
	}

	fileInfo, err := os.Stat(script.AbsolutePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}
	if err := os.WriteFile(targetAbsolutePath, content, fileInfo.Mode().Perm()); err != nil {
		return LocalSkillDetail{}, err
	}

	now := time.Now()
	metadata := skillMetadata{
		CreatedAt:   now,
		UpdatedAt:   now,
		NextVersion: 1,
		Versions:    nil,
	}
	metadata, err = appendSkillVersion(targetRelativePath, metadata, content, now, fmt.Sprintf("Copied from %s", relativePath))
	if err != nil {
		return LocalSkillDetail{}, err
	}
	if err := writeSkillMetadata(targetRelativePath, metadata); err != nil {
		return LocalSkillDetail{}, err
	}

	return getLocalSkillDetail(targetRelativePath)
}

func deleteLocalSkill(relativePath string) error {
	script, err := findManagedScript(relativePath, 0)
	if err != nil {
		return err
	}

	if err := os.Remove(script.AbsolutePath); err != nil {
		return err
	}

	metadataPath := getSkillMetadataPath(script.RelativePath)
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	versionsDir := getSkillVersionsDir(script.RelativePath)
	if err := os.RemoveAll(versionsDir); err != nil {
		return err
	}

	return nil
}

func restoreLocalSkillVersion(relativePath string, versionNumber int) (LocalSkillDetail, error) {
	script, err := findManagedScript(relativePath, defaultPreviewLines)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	metadata, err := ensureSkillMetadata(script)
	if err != nil {
		return LocalSkillDetail{}, err
	}

	var versionFileName string
	for _, version := range metadata.Versions {
		if version.Number == versionNumber {
			versionFileName = version.FileName
			break
		}
	}
	if versionFileName == "" {
		return LocalSkillDetail{}, fmt.Errorf("version v%d was not found", versionNumber)
	}

	versionContent, err := os.ReadFile(filepath.Join(getSkillVersionsDir(relativePath), versionFileName))
	if err != nil {
		return LocalSkillDetail{}, err
	}

	fileInfo, err := os.Stat(script.AbsolutePath)
	if err != nil {
		return LocalSkillDetail{}, err
	}
	if err := os.WriteFile(script.AbsolutePath, versionContent, fileInfo.Mode().Perm()); err != nil {
		return LocalSkillDetail{}, err
	}

	now := time.Now()
	metadata, err = appendSkillVersion(relativePath, metadata, versionContent, now, fmt.Sprintf("Restored from v%d", versionNumber))
	if err != nil {
		return LocalSkillDetail{}, err
	}
	if err := writeSkillMetadata(relativePath, metadata); err != nil {
		return LocalSkillDetail{}, err
	}

	return getLocalSkillDetail(relativePath)
}

func openLocalSkillDir(relativePath string) error {
	script, err := findManagedScript(relativePath, 0)
	if err != nil {
		return err
	}
	return openSystemPath(filepath.Dir(script.AbsolutePath))
}

func buildLocalSkill(script ScriptInfo) (LocalSkill, error) {
	skill, _, err := buildLocalSkillWithMetadata(script)
	return skill, err
}

func buildLocalSkillWithMetadata(script ScriptInfo) (LocalSkill, skillMetadata, error) {
	metadata, err := ensureSkillMetadata(script)
	if err != nil {
		return LocalSkill{}, skillMetadata{}, err
	}

	summary, err := extractSkillSummary(script.AbsolutePath)
	if err != nil {
		return LocalSkill{}, skillMetadata{}, err
	}

	return LocalSkill{
		ID:           script.RelativePath,
		Title:        trimFileExt(script.Name),
		Summary:      summary,
		RelativePath: script.RelativePath,
		AbsolutePath: script.AbsolutePath,
		CreatedAt:    metadata.CreatedAt,
		UpdatedAt:    metadata.UpdatedAt,
		VersionCount: len(metadata.Versions),
		Versions:     metadata.Versions,
	}, metadata, nil
}

func ensureSkillMetadata(script ScriptInfo) (skillMetadata, error) {
	metadata, exists, err := readSkillMetadata(script.RelativePath)
	if err != nil {
		return skillMetadata{}, err
	}

	currentContent, err := os.ReadFile(script.AbsolutePath)
	if err != nil {
		return skillMetadata{}, err
	}

	if !exists {
		metadata = skillMetadata{
			CreatedAt:   script.ModifiedAt,
			UpdatedAt:   script.ModifiedAt,
			NextVersion: 1,
			Versions:    nil,
		}
		metadata, err = appendSkillVersion(script.RelativePath, metadata, currentContent, script.ModifiedAt, "Imported existing script")
		if err != nil {
			return skillMetadata{}, err
		}
		return metadata, writeSkillMetadata(script.RelativePath, metadata)
	}

	if len(metadata.Versions) == 0 {
		metadata.NextVersion = 1
		metadata, err = appendSkillVersion(script.RelativePath, metadata, currentContent, script.ModifiedAt, "Imported existing script")
		if err != nil {
			return skillMetadata{}, err
		}
		return metadata, writeSkillMetadata(script.RelativePath, metadata)
	}

	latestVersion := metadata.Versions[len(metadata.Versions)-1]
	latestContent, err := os.ReadFile(filepath.Join(getSkillVersionsDir(script.RelativePath), latestVersion.FileName))
	if err != nil && !os.IsNotExist(err) {
		return skillMetadata{}, err
	}

	if !bytes.Equal(currentContent, latestContent) {
		metadata, err = appendSkillVersion(script.RelativePath, metadata, currentContent, script.ModifiedAt, "Detected external file change")
		if err != nil {
			return skillMetadata{}, err
		}
		return metadata, writeSkillMetadata(script.RelativePath, metadata)
	}

	if script.ModifiedAt.After(metadata.UpdatedAt) {
		metadata.UpdatedAt = script.ModifiedAt
		if err := writeSkillMetadata(script.RelativePath, metadata); err != nil {
			return skillMetadata{}, err
		}
	}

	return metadata, nil
}

func appendSkillVersion(relativePath string, metadata skillMetadata, content []byte, savedAt time.Time, note string) (skillMetadata, error) {
	if err := ensureSkillStateDirs(relativePath); err != nil {
		return skillMetadata{}, err
	}

	versionNumber := metadata.NextVersion
	if versionNumber <= 0 {
		versionNumber = 1
	}

	extension := filepath.Ext(relativePath)
	fileName := fmt.Sprintf("v%04d%s", versionNumber, extension)
	versionPath := filepath.Join(getSkillVersionsDir(relativePath), fileName)
	if err := os.WriteFile(versionPath, content, 0644); err != nil {
		return skillMetadata{}, err
	}

	metadata.Versions = append(metadata.Versions, SkillVersion{
		Number:   versionNumber,
		Label:    fmt.Sprintf("v%d", versionNumber),
		SavedAt:  savedAt,
		Note:     note,
		FileName: fileName,
	})
	metadata.NextVersion = versionNumber + 1
	if metadata.CreatedAt.IsZero() {
		metadata.CreatedAt = savedAt
	}
	metadata.UpdatedAt = savedAt

	return metadata, nil
}

func readSkillMetadata(relativePath string) (skillMetadata, bool, error) {
	content, err := os.ReadFile(getSkillMetadataPath(relativePath))
	if err != nil {
		if os.IsNotExist(err) {
			return skillMetadata{}, false, nil
		}
		return skillMetadata{}, false, err
	}

	var metadata skillMetadata
	if err := json.Unmarshal(content, &metadata); err != nil {
		return skillMetadata{}, false, err
	}

	return metadata, true, nil
}

func writeSkillMetadata(relativePath string, metadata skillMetadata) error {
	if err := ensureSkillStateDirs(relativePath); err != nil {
		return err
	}

	content, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(getSkillMetadataPath(relativePath), content, 0644)
}

func extractSkillSummary(path string) (string, error) {
	preview, _, err := readPreview(path, 12)
	if err != nil {
		return "", err
	}

	for _, line := range preview {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#!") || trimmed == "<?php" {
			continue
		}

		trimmed = strings.TrimSuffix(trimmed, "*/")
		for _, prefix := range []string{"#", "//", "--", ";", "/*", "*"} {
			if strings.HasPrefix(trimmed, prefix) {
				trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				break
			}
		}

		if trimmed != "" {
			return trimmed, nil
		}
	}

	return "No description yet", nil
}

func trimFileExt(name string) string {
	extension := filepath.Ext(name)
	return strings.TrimSuffix(name, extension)
}

func normalizeSkillRelativePath(relativePath string) (string, error) {
	cleaned, err := sanitizeSkillRelativePath(relativePath)
	if err != nil {
		return "", err
	}
	if filepath.Ext(cleaned) == "" {
		cleaned += ".sh"
	}

	return cleaned, nil
}

func normalizeSkillRelativePathWithRequiredExtension(relativePath string) (string, error) {
	cleaned, err := sanitizeSkillRelativePath(relativePath)
	if err != nil {
		return "", err
	}
	if filepath.Ext(cleaned) == "" {
		return "", fmt.Errorf("managed skill path must include a file extension")
	}

	return cleaned, nil
}

func sanitizeSkillRelativePath(relativePath string) (string, error) {
	trimmed := strings.TrimSpace(relativePath)
	if trimmed == "" {
		return "", fmt.Errorf("managed skill path cannot be empty")
	}
	if filepath.IsAbs(trimmed) {
		return "", fmt.Errorf("managed skill path must be relative")
	}

	cleaned := filepath.ToSlash(filepath.Clean(filepath.FromSlash(trimmed)))
	if cleaned == "." || cleaned == "" {
		return "", fmt.Errorf("managed skill path cannot be empty")
	}
	if strings.HasPrefix(cleaned, "../") || cleaned == ".." {
		return "", fmt.Errorf("managed skill path cannot escape the scripts directory")
	}

	return cleaned, nil
}

func defaultSkillTemplate(relativePath string) string {
	title := trimFileExt(filepath.Base(relativePath))
	extension := strings.ToLower(filepath.Ext(relativePath))

	switch extension {
	case ".sh":
		return fmt.Sprintf("#!/usr/bin/env sh\n# %s\n\n", title)
	case ".bash":
		return fmt.Sprintf("#!/usr/bin/env bash\n# %s\n\n", title)
	case ".zsh":
		return fmt.Sprintf("#!/usr/bin/env zsh\n# %s\n\n", title)
	case ".py":
		return fmt.Sprintf("#!/usr/bin/env python3\n# %s\n\n", title)
	case ".js", ".cjs", ".mjs":
		return fmt.Sprintf("#!/usr/bin/env node\n// %s\n\n", title)
	case ".rb":
		return fmt.Sprintf("#!/usr/bin/env ruby\n# %s\n\n", title)
	case ".pl":
		return fmt.Sprintf("#!/usr/bin/env perl\n# %s\n\n", title)
	case ".lua":
		return fmt.Sprintf("#!/usr/bin/env lua\n-- %s\n\n", title)
	case ".php":
		return fmt.Sprintf("#!/usr/bin/env php\n<?php\n\n// %s\n\n", title)
	case ".ps1":
		return fmt.Sprintf("# %s\n\n", title)
	case ".ts", ".tsx", ".jsx":
		return fmt.Sprintf("// %s\n\n", title)
	default:
		return fmt.Sprintf("# %s\n\n", title)
	}
}

func nextCopyRelativePath(relativePath string) (string, error) {
	dir := filepath.Dir(relativePath)
	base := filepath.Base(relativePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	candidates := []string{
		fmt.Sprintf("%s copy%s", name, ext),
	}
	for i := 2; i < 1000; i++ {
		candidates = append(candidates, fmt.Sprintf("%s copy %d%s", name, i, ext))
	}

	for _, candidate := range candidates {
		nextPath := candidate
		if dir != "." {
			nextPath = filepath.ToSlash(filepath.Join(dir, candidate))
		}
		if _, err := os.Stat(filepath.Join(getScriptsDir(), filepath.FromSlash(nextPath))); os.IsNotExist(err) {
			return nextPath, nil
		}
	}

	return "", fmt.Errorf("failed to generate a copy name for %s", relativePath)
}
