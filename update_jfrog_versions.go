package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
)

func main() {
	fmt.Println("🚀 JFrog CLI Artifactory Version Updater")
	fmt.Println("==========================================")

	// JFrog repositories to update
	jfrogRepos := []string{
		"github.com/jfrog/build-info-go",
		"github.com/jfrog/froggit-go",
		"github.com/jfrog/gofrog",
		"github.com/jfrog/jfrog-cli-core/v2",
		"github.com/jfrog/jfrog-client-go",
	}

	// Get current directory and construct paths
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	// Path to jfrog-cli-artifactory go.mod
	// Check if we're in automation-scripts directory
	var goModPath, jfrogCliArtifactoryDir string
	if filepath.Base(currentDir) == "automation-scripts" {
		// We're in automation-scripts, look in parent directory
		parentDir := filepath.Dir(currentDir)
		goModPath = filepath.Join(parentDir, "jfrog-cli-artifactory", "go.mod")
		jfrogCliArtifactoryDir = filepath.Join(parentDir, "jfrog-cli-artifactory")
	} else {
		// We're in root directory
		goModPath = filepath.Join(currentDir, "jfrog-cli-artifactory", "go.mod")
		jfrogCliArtifactoryDir = filepath.Join(currentDir, "jfrog-cli-artifactory")
	}

	// Check if jfrog-cli-artifactory directory exists
	if _, err := os.Stat(jfrogCliArtifactoryDir); os.IsNotExist(err) {
		fmt.Printf("❌ jfrog-cli-artifactory directory not found at: %s\n", jfrogCliArtifactoryDir)
		fmt.Println("💡 Make sure to run this script from:")
		fmt.Println("   - The root directory containing jfrog-cli-artifactory")
		fmt.Println("   - OR from the automation-scripts directory")
		os.Exit(1)
	}

	fmt.Printf("📁 Target: %s\n\n", goModPath)

	// Read go.mod file
	data, err := os.ReadFile(goModPath)
	if err != nil {
		fmt.Printf("❌ Failed to read %s: %v\n", goModPath, err)
		os.Exit(1)
	}

	// Parse go.mod
	f, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		fmt.Printf("❌ Failed to parse go.mod: %v\n", err)
		os.Exit(1)
	}

	// Update each JFrog repository
	updated := 0
	for _, repo := range jfrogRepos {
		fmt.Printf("🔍 Checking %s...\n", repo)

		latestVersion, err := getLatestVersion(repo)
		if err != nil {
			fmt.Printf("⚠️  Failed to get latest version for %s: %v\n", repo, err)
			continue
		}

		// Get current version for comparison
		currentVersion := getCurrentVersion(f, repo)
		if currentVersion == latestVersion {
			fmt.Printf("✨ %s is already up to date (%s)\n", repo, latestVersion)
			continue
		}

		// Update require directive
		if err := f.AddRequire(repo, latestVersion); err != nil {
			fmt.Printf("⚠️  Failed to update require for %s: %v\n", repo, err)
			continue
		}

		// Update replace directive if exists
		for _, r := range f.Replace {
			if r.Old.Path == repo {
				if err := f.DropReplace(repo, ""); err != nil {
					fmt.Printf("⚠️  Failed to drop replace for %s: %v\n", repo, err)
					continue
				}
				if err := f.AddReplace(repo, "", repo, latestVersion); err != nil {
					fmt.Printf("⚠️  Failed to add replace for %s: %v\n", repo, err)
					continue
				}
				break
			}
		}

		if currentVersion != "" {
			fmt.Printf("✅ Updated %s: %s → %s\n", repo, currentVersion, latestVersion)
		} else {
			fmt.Printf("✅ Updated %s to %s\n", repo, latestVersion)
		}
		updated++
	}

	if updated == 0 {
		fmt.Println("\n🎯 All JFrog repositories are already up to date!")
		return
	}

	// Format the go.mod file
	f.Cleanup()

	// Write back to file
	formatted, err := f.Format()
	if err != nil {
		fmt.Printf("❌ Failed to format go.mod: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(goModPath, formatted, 0644); err != nil {
		fmt.Printf("❌ Failed to write go.mod: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n🎉 Successfully updated %d JFrog repositories!\n", updated)

	// Automatically run go mod tidy
	fmt.Println("🧹 Running 'go mod tidy'...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = jfrogCliArtifactoryDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Failed to run 'go mod tidy': %v\n", err)
		fmt.Println("💡 Please run 'go mod tidy' manually in jfrog-cli-artifactory directory")
	} else {
		fmt.Println("✅ Successfully ran 'go mod tidy'")
	}

	fmt.Println("\n🎊 All done! JFrog repository versions updated and dependencies cleaned up.")
}

func getCurrentVersion(f *modfile.File, repo string) string {
	for _, req := range f.Require {
		if req.Mod.Path == repo {
			return req.Mod.Version
		}
	}
	return ""
}

func getLatestVersion(repo string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Get latest release from GitHub API
	parts := strings.Split(repo, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid repo format: %s", repo)
	}

	owner := parts[1]
	repoName := parts[2]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repoName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Add GitHub token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	// Ensure version has 'v' prefix as required by Go modules
	version := release.TagName
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}
	return version, nil
}
