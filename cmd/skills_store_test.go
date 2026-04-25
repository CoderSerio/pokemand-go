package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillLifecycle(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("PKMG_CONFIG_DIR", filepath.Join(tempDir, "config"))
	t.Setenv("PKMG_DATA_DIR", filepath.Join(tempDir, "data"))
	resetPathCaches()
	defer resetPathCaches()

	created, err := createLocalSkill("team/deploy", "#!/bin/sh\n# Deploy skill\necho deploy\n")
	if err != nil {
		t.Fatalf("createLocalSkill failed: %v", err)
	}
	if created.RelativePath != "team/deploy.sh" {
		t.Fatalf("unexpected relative path: %s", created.RelativePath)
	}
	if created.VersionCount != 1 {
		t.Fatalf("expected version count 1, got %d", created.VersionCount)
	}

	results, err := searchLocalSkills("deploy")
	if err != nil {
		t.Fatalf("searchLocalSkills failed: %v", err)
	}
	if len(results) != 1 || results[0].RelativePath != "team/deploy.sh" {
		t.Fatalf("unexpected search results: %+v", results)
	}

	updated, err := saveLocalSkill("team/deploy.sh", "#!/bin/sh\n# Deploy skill updated\necho deploy-now\n")
	if err != nil {
		t.Fatalf("saveLocalSkill failed: %v", err)
	}
	if updated.VersionCount != 2 {
		t.Fatalf("expected version count 2 after save, got %d", updated.VersionCount)
	}
	if updated.Summary != "Deploy skill updated" {
		t.Fatalf("unexpected summary after save: %s", updated.Summary)
	}

	detail, err := getLocalSkillDetail("team/deploy.sh")
	if err != nil {
		t.Fatalf("getLocalSkillDetail failed: %v", err)
	}
	if detail.Content != "#!/bin/sh\n# Deploy skill updated\necho deploy-now\n" {
		t.Fatalf("unexpected content after save: %q", detail.Content)
	}

	if err := deleteLocalSkill("team/deploy.sh"); err != nil {
		t.Fatalf("deleteLocalSkill failed: %v", err)
	}

	if _, err := findManagedScript("team/deploy.sh", 0); !os.IsNotExist(err) {
		t.Fatalf("expected deleted skill to be missing, got err=%v", err)
	}

	results, err = searchLocalSkills("deploy")
	if err != nil {
		t.Fatalf("searchLocalSkills after delete failed: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected no search results after delete, got %+v", results)
	}

	metadataPath := getSkillMetadataPath("team/deploy.sh")
	if _, err := os.Stat(metadataPath); !os.IsNotExist(err) {
		t.Fatalf("expected metadata to be removed, stat err=%v", err)
	}
}
