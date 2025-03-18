package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	repoOwner = "pagerguild"
	repoName  = "guilde-cli-releases"
)

const (
	assetNameMacosArm64 = "guilde-cli-darwin-arm64"
	assetNameMacosIntel = "guilde-cli-darwin-amd64"
	assetNameLinuxArm64 = "guilde-cli-linux-arm64"
	assetNameLinuxIntel = "guilde-cli-linux-amd64"
)

const (
	releaseNotesFileName = "RELEASE_NOTES.md"
)

const (
	// We use zip files, since these are supported by Apple Developer's
	// code signing tools.
	fileType = "zip"
)

// ReleaseAsset represents a binary release asset for a specific platform/arch
type ReleaseAsset struct {
	Name     string // Base name without extension (e.g., "guilde-cli-darwin-arm64")
	Path     string // Full path to the asset file
	Checksum string // SHA256 checksum (computed)
}

// Release represents all the information needed for a release
type Release struct {
	Version      string
	Assets       []ReleaseAsset
	NotesPath    string
	NotesContent string
}

// NewRelease creates and validates a new Release from a directory
func NewRelease(version string, dirPath string) (*Release, error) {
	assetNames := []string{
		assetNameMacosArm64,
		assetNameMacosIntel,
		assetNameLinuxArm64,
		assetNameLinuxIntel,
	}

	assets := make([]ReleaseAsset, 0, len(assetNames))
	for _, name := range assetNames {
		asset := ReleaseAsset{
			Name: name,
			Path: filepath.Join(dirPath, name+"."+fileType),
		}
		assets = append(assets, asset)
	}

	notesPath := filepath.Join(dirPath, releaseNotesFileName)

	r := &Release{
		Version:   version,
		Assets:    assets,
		NotesPath: notesPath,
	}

	if err := r.Validate(); err != nil {
		return nil, fmt.Errorf("release validation failed: %w", err)
	}

	if err := r.ComputeChecksums(); err != nil {
		return nil, fmt.Errorf("checksum computation failed: %w", err)
	}

	if err := r.LoadReleaseNotes(); err != nil {
		return nil, fmt.Errorf("failed to load release notes: %w", err)
	}

	return r, nil
}

func (r *Release) Validate() error {
	// Check release notes exist
	if _, err := os.Stat(r.NotesPath); os.IsNotExist(err) {
		return fmt.Errorf("release notes not found at %s", r.NotesPath)
	}

	// Check all assets exist
	for _, asset := range r.Assets {
		if _, err := os.Stat(asset.Path); os.IsNotExist(err) {
			return fmt.Errorf("asset not found at %s", asset.Path)
		}
	}

	return nil
}

func (r *Release) ComputeChecksums() error {
	for i := range r.Assets {
		checksum, err := calculateSHA256(r.Assets[i].Path)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum for %s: %w", r.Assets[i].Name, err)
		}
		r.Assets[i].Checksum = checksum
	}
	return nil
}

func (r *Release) LoadReleaseNotes() error {
	content, err := os.ReadFile(r.NotesPath)
	if err != nil {
		return fmt.Errorf("failed to read release notes: %w", err)
	}
	r.NotesContent = string(content)
	return nil
}

func calculateSHA256(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s RELEASE_VERSION DIRECTORY_PATH\n", os.Args[0])
		os.Exit(1)
	}

	version := os.Args[1]
	dirPath := os.Args[2]

	// Validate version format (basic check)
	if len(version) == 0 || version[0] == 'v' {
		fmt.Fprintf(os.Stderr, "Version should be in format X.Y.Z (without 'v' prefix)\n")
		os.Exit(1)
	}

	// Check for GITHUB_PAT
	if os.Getenv("GITHUB_PAT") == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_PAT environment variable is required\n")
		os.Exit(1)
	}

	release, err := NewRelease(version, dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create release: %v\n", err)
		os.Exit(1)
	}

	// Print checksums (for now)
	for _, asset := range release.Assets {
		fmt.Printf("SHA256 (%s.%s) = %s\n", asset.Name, fileType, asset.Checksum)
	}

	// TODO: Create GitHub release using the PAT
	// This will be implemented in the next step
}
