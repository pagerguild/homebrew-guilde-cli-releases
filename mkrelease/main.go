package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//go:generate mockgen -source=main.go -destination=./mock_test.go -package=main

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

// Release defines the interface for release operations
type Release interface {
	Validate() error
	ComputeChecksums() error
	LoadReleaseNotes() error
	CreateAssets(dirPath string) error
	GetVersion() string
	GetAssets() []ReleaseAsset
	GetNotesContent() string
}

// ReleaseImpl implements the Release interface
type ReleaseImpl struct {
	version      string
	assets       []ReleaseAsset
	notesPath    string
	notesContent string
}

// GetVersion returns the release version
func (r *ReleaseImpl) GetVersion() string {
	return r.version
}

// GetAssets returns the release assets
func (r *ReleaseImpl) GetAssets() []ReleaseAsset {
	return r.assets
}

// GetNotesContent returns the release notes content
func (r *ReleaseImpl) GetNotesContent() string {
	return r.notesContent
}

// CreateAssets populates the release assets from the given directory
func (r *ReleaseImpl) CreateAssets(dirPath string) error {
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

	r.assets = assets
	r.notesPath = filepath.Join(dirPath, releaseNotesFileName)
	return nil
}

func (r *ReleaseImpl) Validate() error {
	// Check release notes exist
	if _, err := os.Stat(r.notesPath); os.IsNotExist(err) {
		return fmt.Errorf("release notes not found at %s", r.notesPath)
	}

	// Check all assets exist
	for _, asset := range r.assets {
		if _, err := os.Stat(asset.Path); os.IsNotExist(err) {
			return fmt.Errorf("asset not found at %s", asset.Path)
		}
	}

	return nil
}

func (r *ReleaseImpl) ComputeChecksums() error {
	for i := range r.assets {
		checksum, err := calculateSHA256(r.assets[i].Path)
		if err != nil {
			return fmt.Errorf("failed to calculate checksum for %s: %w", r.assets[i].Name, err)
		}
		r.assets[i].Checksum = checksum
	}
	return nil
}

func (r *ReleaseImpl) LoadReleaseNotes() error {
	content, err := os.ReadFile(r.notesPath)
	if err != nil {
		return fmt.Errorf("failed to read release notes: %w", err)
	}
	r.notesContent = string(content)
	return nil
}

// ReleaseStrategy defines the steps to create a release
func ReleaseStrategy(r Release, dirPath string) error {
	if err := r.CreateAssets(dirPath); err != nil {
		return fmt.Errorf("failed to create assets: %w", err)
	}

	if err := r.Validate(); err != nil {
		return fmt.Errorf("release validation failed: %w", err)
	}

	if err := r.ComputeChecksums(); err != nil {
		return fmt.Errorf("checksum computation failed: %w", err)
	}

	if err := r.LoadReleaseNotes(); err != nil {
		return fmt.Errorf("failed to load release notes: %w", err)
	}

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

	// Strip 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	// Validate version is not empty
	if len(version) == 0 {
		fmt.Fprintf(os.Stderr, "Version cannot be empty\n")
		os.Exit(1)
	}

	// Check for GITHUB_PAT
	if os.Getenv("GITHUB_PAT") == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_PAT environment variable is required\n")
		os.Exit(1)
	}

	release := &ReleaseImpl{version: version}
	if err := ReleaseStrategy(release, dirPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create release: %v\n", err)
		os.Exit(1)
	}

	// Print checksums (for now)
	for _, asset := range release.GetAssets() {
		fmt.Printf("SHA256 (%s.%s) = %s\n", asset.Name, fileType, asset.Checksum)
	}

	// TODO: Create GitHub release using the PAT
	// This will be implemented in the next step
}
