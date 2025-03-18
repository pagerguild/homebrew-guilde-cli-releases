package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v69/github"
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
	fileType              = "zip"
	mediaType             = "application/zip"
	releaseNotesMediaType = "text/markdown"
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
	GetRepoPath() string
	// Uses go-get to create a new tag called `vX.X.X` where X.X.X is the version.
	// If the tag already exists continue; just push the tags to origin.
	CreateAndPushTagToGitHub(ctx context.Context, pat, version, path string) error
	CreateGitHubRelease(ctx context.Context, client *github.Client, version, notesContent string) (int64, error)
	UploadReleaseAsset(ctx context.Context, client *github.Client, releaseID int64, name, mediaType string, file *os.File) error
	ReleaseExists(ctx context.Context, client *github.Client, version string) (bool, error)
	// CreateGitHubClient creates a new GitHub client that includes the
	// authentication token in each request.
	CreateGitHubClient(pat string) *github.Client
}

// ReleaseImpl implements the Release interface
type ReleaseImpl struct {
	version      string
	assets       []ReleaseAsset
	notesPath    string
	notesContent string
	repoPath     string
}

func tagName(version string) string {
	return "v" + version
}

// ReleaseExists implements Release.
func (r *ReleaseImpl) ReleaseExists(ctx context.Context, client *github.Client, version string) (bool, error) {
	_, resp, err := client.Repositories.GetReleaseByTag(ctx, repoOwner, repoName, tagName(version))
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get release: %w", err)
	}
	return true, nil
}

// CreateAndPushTagToGitHub implements Release.
func (r *ReleaseImpl) CreateAndPushTagToGitHub(ctx context.Context, pat string, version string, path string) error {
	// Open the repository
	repo, err := git.PlainOpen(r.repoPath)
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get the HEAD reference
	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Create tag name (v + version)
	tagName := "v" + version

	// Check if the tag already exists
	_, err = repo.Tag(tagName)
	if err == nil {
		// Tag already exists, skip creation but still need to push
		fmt.Printf("Tag %s already exists, skipping creation\n", tagName)
	} else {
		// Create a new tag
		_, err = repo.CreateTag(tagName, headRef.Hash(), &git.CreateTagOptions{
			Message: fmt.Sprintf("Release %s", tagName),
		})
		if err != nil {
			return fmt.Errorf("failed to create tag: %w", err)
		}
		fmt.Printf("Created tag %s\n", tagName)
	}

	// Get the remote
	remote, err := repo.Remote("origin")
	if err != nil {
		return fmt.Errorf("failed to get origin remote: %w", err)
	}

	// Push the tag to origin
	refspec := fmt.Sprintf("refs/tags/%s:refs/tags/%s", tagName, tagName)
	err = remote.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(refspec)},
		Auth: &githttp.BasicAuth{
			Username: "git", // This can be any non-empty string
			Password: pat,   // PAT is used as the password
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("failed to push tag: %w", err)
	}

	fmt.Printf("Pushed tag %s to origin\n", tagName)
	return nil
}

// CreateGitHubClient implements Release.
func (r *ReleaseImpl) CreateGitHubClient(pat string) *github.Client {
	// Create an authenticated client with the provided personal access token
	client := github.NewClient(nil).WithAuthToken(pat)
	return client
}

// CreateGitHubRelease implements Release.
func (r *ReleaseImpl) CreateGitHubRelease(ctx context.Context, client *github.Client, version string, notesContent string) (int64, error) {
	// Check if release already exists
	exists, err := r.ReleaseExists(ctx, client, version)
	if err != nil {
		return 0, fmt.Errorf("failed to check if release exists: %w", err)
	}
	if exists {
		return 0, fmt.Errorf("release %s already exists", tagName(version))
	}

	// Create a new release
	tagName := tagName(version)
	release, _, err := client.Repositories.CreateRelease(ctx, repoOwner, repoName, &github.RepositoryRelease{
		TagName:         &tagName,
		Name:            &tagName,
		Body:            &notesContent,
		TargetCommitish: nil, // Use the default branch (main)
		Draft:           github.Ptr(false),
		Prerelease:      github.Ptr(false),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create release: %w", err)
	}

	fmt.Printf("Created GitHub release %s with ID %d\n", tagName, release.GetID())
	return release.GetID(), nil
}

// UploadReleaseAsset implements Release.
func (r *ReleaseImpl) UploadReleaseAsset(ctx context.Context, client *github.Client, releaseID int64, name string, mediaType string, file *os.File) error {
	// Get file info for size calculation
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Prepare the asset options
	opts := &github.UploadOptions{
		Name:      name + "." + fileType,
		MediaType: mediaType,
	}

	// Upload the asset to the release
	_, _, err = client.Repositories.UploadReleaseAsset(ctx, repoOwner, repoName, releaseID, opts, file)
	if err != nil {
		return fmt.Errorf("failed to upload release asset: %w", err)
	}

	fmt.Printf("Uploaded asset %s (%.2f KB) to release %d\n", opts.Name, float64(fileInfo.Size())/1024.0, releaseID)
	return nil
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

// GetRepoPath returns the repository path
func (r *ReleaseImpl) GetRepoPath() string {
	return r.repoPath
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

	// Check if the repo path is a valid git repository
	if r.repoPath != "" {
		_, err := git.PlainOpen(r.repoPath)
		if err != nil {
			return fmt.Errorf("failed to open git repository at %s: %w", r.repoPath, err)
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
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s RELEASE_VERSION DIRECTORY_PATH REPO_PATH\n", os.Args[0])
		os.Exit(1)
	}

	version := os.Args[1]
	dirPath := os.Args[2]
	repoPath := os.Args[3]

	// Strip 'v' prefix if present
	version = strings.TrimPrefix(version, "v")

	// Validate version is not empty
	if len(version) == 0 {
		fmt.Fprintf(os.Stderr, "Version cannot be empty\n")
		os.Exit(1)
	}

	// Check for GITHUB_PAT
	pat := os.Getenv("GITHUB_PAT")
	if pat == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_PAT environment variable is required\n")
		os.Exit(1)
	}

	release := &ReleaseImpl{
		version:  version,
		repoPath: repoPath,
	}

	// Step 1: Set up assets and validate
	if err := ReleaseStrategy(release, dirPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create release: %v\n", err)
		os.Exit(1)
	}

	// Print checksums
	for _, asset := range release.GetAssets() {
		fmt.Printf("SHA256 (%s.%s) = %s\n", asset.Name, fileType, asset.Checksum)
	}

	// Create a context for GitHub operations
	ctx := context.Background()

	// Step 2: Create and push tag
	if err := release.CreateAndPushTagToGitHub(ctx, pat, version, repoPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create and push tag: %v\n", err)
		os.Exit(1)
	}

	// Step 3: Create GitHub client
	client := release.CreateGitHubClient(pat)

	// Step 4: Create GitHub release
	releaseID, err := release.CreateGitHubRelease(ctx, client, version, release.GetNotesContent())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create GitHub release: %v\n", err)
		os.Exit(1)
	}

	// Step 5: Upload assets
	for _, asset := range release.GetAssets() {
		// Open the asset file
		file, err := os.Open(asset.Path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open asset file %s: %v\n", asset.Path, err)
			os.Exit(1)
		}

		// Upload the asset
		err = release.UploadReleaseAsset(ctx, client, releaseID, asset.Name, mediaType, file)
		if err != nil {
			file.Close()
			fmt.Fprintf(os.Stderr, "Failed to upload asset %s: %v\n", asset.Name, err)
			os.Exit(1)
		}
		file.Close()
	}

	fmt.Printf("Successfully created release v%s with all assets\n", version)
}
