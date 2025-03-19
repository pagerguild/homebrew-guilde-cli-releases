package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"text/template"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v69/github"
)

//go:generate mockgen -source=main.go -destination=./mock_test.go -package=main

const (
	repoOwner = "pagerguild"
	repoName  = "guilde-cli-releases"
)

const (
	sourceRepoOwner = "pagerguild"
	sourceRepoName  = "pagerguild"
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

	// Formula constants
	formulaFileName = "guilde-cli.rb"
)

//go:embed guilde-cli.rb.tmpl
var formulaTemplate string

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
	// Opens an asset file for reading. This makes file operations mockable.
	OpenAssetFile(path string) (*os.File, error)
	// Creates and pushes a tag to the specified commit
	CreateAndPushTagToGitHub(ctx context.Context, pat, version, path, commitSHA string) error
	CreateGitHubRelease(ctx context.Context, client *github.Client, version, notesContent string) (int64, error)
	UploadReleaseAsset(ctx context.Context, client *github.Client, releaseID int64, name, mediaType string, file *os.File) error
	ReleaseExists(ctx context.Context, client *github.Client, version string) (bool, error)
	// CreateGitHubClient creates a new GitHub client that includes the
	// authentication token in each request.
	CreateGitHubClient(pat string) *github.Client
	// Gets the formula template as a string
	GetFormulaTemplate() string
	// Gets the path where the formula file should be saved
	GetFormulaFilePath() string
	// Renders the formula template with the current release data
	RenderFormulaTemplate() (string, error)
	// Saves the rendered formula to the formula file
	SaveFormulaFile(content string) error
	// Commits the formula change and returns the commit SHA
	CommitFormulaChange(ctx context.Context, pat, message string) (string, error)
	// Downloads release assets and notes from the pagerguild/pagerguild repository
	DownloadReleaseAssetsAndNotes(ctx context.Context, client *github.Client, dirPath string, pat string) error
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
func (r *ReleaseImpl) CreateAndPushTagToGitHub(ctx context.Context, pat string, version string, path string, commitSHA string) error {
	// Open the repository
	repo, err := git.PlainOpen(r.repoPath)
	if err != nil {
		return fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get the reference to the specified commit or HEAD if not provided
	var hash plumbing.Hash
	if commitSHA != "" {
		hash = plumbing.NewHash(commitSHA)
	} else {
		// Get the HEAD reference
		headRef, err := repo.Head()
		if err != nil {
			return fmt.Errorf("failed to get HEAD reference: %w", err)
		}
		hash = headRef.Hash()
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
		_, err = repo.CreateTag(tagName, hash, &git.CreateTagOptions{
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

// OpenAssetFile implements the file opening functionality from the Release interface
func (r *ReleaseImpl) OpenAssetFile(path string) (*os.File, error) {
	return os.Open(path)
}

// DownloadReleaseAssetsAndNotes downloads release assets and notes from the pagerguild/pagerguild repository
func (r *ReleaseImpl) DownloadReleaseAssetsAndNotes(ctx context.Context, client *github.Client, dirPath, pat string) error {
	version := r.GetVersion()
	tag := tagName(version)

	// Get the release by tag
	release, _, err := client.Repositories.GetReleaseByTag(ctx, sourceRepoOwner, sourceRepoName, tag)
	if err != nil {
		return fmt.Errorf("failed to get release from %s/%s with tag %s: %w", sourceRepoOwner, sourceRepoName, tag, err)
	}

	// Get release notes
	r.notesContent = release.GetBody()
	r.notesPath = filepath.Join(dirPath, releaseNotesFileName)

	// Save release notes to file
	err = os.WriteFile(r.notesPath, []byte(r.notesContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to save release notes: %w", err)
	}
	fmt.Printf("Saved release notes to %s\n", r.notesPath)

	// Prepare assets
	assetNames := []string{
		assetNameMacosArm64,
		assetNameMacosIntel,
		assetNameLinuxArm64,
		assetNameLinuxIntel,
	}
	assets := make([]ReleaseAsset, 0, len(assetNames))

	// Download each asset
	for _, name := range assetNames {
		// Look for the asset in the release
		var asset *github.ReleaseAsset

		fmt.Printf("Looking for asset %s.%s in release with %d assets\n", name, fileType, len(release.Assets))
		for _, a := range release.Assets {
			downloadName := fmt.Sprintf("%s.%s", name, fileType)
			if a.GetName() == downloadName {
				asset = a
				break
			}
		}

		if asset == nil {
			return fmt.Errorf("asset %s.%s not found in release", name, fileType)
		}

		// Create the local file
		localPath := filepath.Join(dirPath, fmt.Sprintf("%s.%s", name, fileType))
		fmt.Printf("Downloading asset %s (ID: %d) to %s\n", name, asset.GetID(), localPath)

		// Create the output file
		out, err := os.Create(localPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", localPath, err)
		}

		// Download the asset using GitHub client
		// Pass nil as followRedirectsClient to get redirectURL instead of io.Reader
		rc, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, sourceRepoOwner, sourceRepoName, asset.GetID(), AuthenticatedHttpClient(pat))
		if err != nil {
			out.Close()
			return fmt.Errorf("failed to get asset download info: %w", err)
		}

		// Per the docs, exactly one of rc and redirectURL will be non-zero
		if rc != nil {
			// Direct download case
			_, err = io.Copy(out, rc)
			rc.Close()
			if err != nil {
				out.Close()
				return fmt.Errorf("failed to save asset to file: %w", err)
			}
		} else if redirectURL != "" {
			// Redirect case
			resp, err := http.Get(redirectURL)
			if err != nil {
				out.Close()
				return fmt.Errorf("failed to download from redirect URL: %w", err)
			}
			defer resp.Body.Close()

			_, err = io.Copy(out, resp.Body)
			if err != nil {
				out.Close()
				return fmt.Errorf("failed to save asset from redirect to file: %w", err)
			}
		}

		out.Close()

		// Add to our assets list
		assetInfo := ReleaseAsset{
			Name: name,
			Path: localPath,
		}
		assets = append(assets, assetInfo)
	}

	r.assets = assets
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

// GetFormulaTemplate returns the Homebrew formula template
func (r *ReleaseImpl) GetFormulaTemplate() string {
	return formulaTemplate
}

// GetFormulaFilePath returns the path to the formula file
func (r *ReleaseImpl) GetFormulaFilePath() string {
	return filepath.Join(r.repoPath, formulaFileName)
}

// RenderFormulaTemplate renders the formula template with release data
func (r *ReleaseImpl) RenderFormulaTemplate() (string, error) {
	// Find the asset checksums
	var macOSArmChecksum, macOSIntelChecksum, linuxArmChecksum, linuxIntelChecksum string
	for _, asset := range r.assets {
		switch asset.Name {
		case assetNameMacosArm64:
			macOSArmChecksum = asset.Checksum
		case assetNameMacosIntel:
			macOSIntelChecksum = asset.Checksum
		case assetNameLinuxArm64:
			linuxArmChecksum = asset.Checksum
		case assetNameLinuxIntel:
			linuxIntelChecksum = asset.Checksum
		}
	}

	// Create the template data
	data := struct {
		Version            string
		MacOSArmChecksum   string
		MacOSIntelChecksum string
		LinuxArmChecksum   string
		LinuxIntelChecksum string
	}{
		Version:            r.version,
		MacOSArmChecksum:   macOSArmChecksum,
		MacOSIntelChecksum: macOSIntelChecksum,
		LinuxArmChecksum:   linuxArmChecksum,
		LinuxIntelChecksum: linuxIntelChecksum,
	}

	// Parse the template
	tmpl, err := template.New("formula").Parse(r.GetFormulaTemplate())
	if err != nil {
		return "", fmt.Errorf("failed to parse formula template: %w", err)
	}

	// Render the template
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("failed to render formula template: %w", err)
	}

	return rendered.String(), nil
}

// SaveFormulaFile saves the rendered formula to a file
func (r *ReleaseImpl) SaveFormulaFile(content string) error {
	formulaPath := r.GetFormulaFilePath()
	err := os.WriteFile(formulaPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write formula file: %w", err)
	}
	return nil
}

// CommitFormulaChange commits the formula change and returns the commit SHA
func (r *ReleaseImpl) CommitFormulaChange(ctx context.Context, pat, message string) (string, error) {
	// Open the repository
	repo, err := git.PlainOpen(r.repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get the worktree
	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	// Add the formula file to the staging area
	formulaPath := r.GetFormulaFilePath()
	_, err = worktree.Add(filepath.Base(formulaPath))
	if err != nil {
		return "", fmt.Errorf("failed to add formula file to staging: %w", err)
	}

	// Commit the change
	commitOptions := &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Guilde CLI Release Bot",
			Email: "bot@pagerguild.com",
			When:  time.Now(),
		},
	}

	hash, err := worktree.Commit(message, commitOptions)
	if err != nil {
		return "", fmt.Errorf("failed to commit formula change: %w", err)
	}

	// Push the commit to the origin
	err = repo.Push(&git.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: "git", // This can be any non-empty string
			Password: pat,   // PAT is used as the password
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return "", fmt.Errorf("failed to push commit: %w", err)
	}

	return hash.String(), nil
}

// ReleaseStrategy defines the steps to create a release
func ReleaseStrategy(ctx context.Context, r Release, dirPath string, pat string) error {
	// Create GitHub client (needed for download and release checks)
	client := r.CreateGitHubClient(pat)

	// Step 1: Download assets and release notes from pagerguild/pagerguild
	if err := r.DownloadReleaseAssetsAndNotes(ctx, client, dirPath, pat); err != nil {
		return fmt.Errorf("failed to download release assets: %w", err)
	}

	// Step 2: Validate assets
	if err := r.Validate(); err != nil {
		return fmt.Errorf("release validation failed: %w", err)
	}

	if err := r.ComputeChecksums(); err != nil {
		return fmt.Errorf("checksum computation failed: %w", err)
	}

	// Print checksums
	for _, asset := range r.GetAssets() {
		fmt.Printf("SHA256 (%s.%s) = %s\n", asset.Name, fileType, asset.Checksum)
	}

	version := r.GetVersion()

	// Step 3: Check if release already exists
	exists, err := r.ReleaseExists(ctx, client, version)
	if err != nil {
		return fmt.Errorf("failed to check if release exists: %w", err)
	}
	if exists {
		return fmt.Errorf("release %s already exists", tagName(version))
	}

	// Step 4: Generate and save the Homebrew formula
	formula, err := r.RenderFormulaTemplate()
	if err != nil {
		return fmt.Errorf("failed to render formula: %w", err)
	}

	if err := r.SaveFormulaFile(formula); err != nil {
		return fmt.Errorf("failed to save formula file: %w", err)
	}

	// Step 5: Commit the formula change
	commitMsg := fmt.Sprintf("Update formula for release v%s", version)
	commitSHA, err := r.CommitFormulaChange(ctx, pat, commitMsg)
	if err != nil {
		return fmt.Errorf("failed to commit formula change: %w", err)
	}

	fmt.Printf("Committed formula update with SHA: %s\n", commitSHA)

	// Step 6: Create and push tag to the new commit
	if err := r.CreateAndPushTagToGitHub(ctx, pat, version, r.GetRepoPath(), commitSHA); err != nil {
		return fmt.Errorf("failed to create and push tag: %w", err)
	}

	// Step 7: Create GitHub release
	releaseID, err := r.CreateGitHubRelease(ctx, client, version, r.GetNotesContent())
	if err != nil {
		return fmt.Errorf("failed to create GitHub release: %w", err)
	}

	// Step 8: Upload assets
	for _, asset := range r.GetAssets() {
		// Open the asset file using the interface method
		file, err := r.OpenAssetFile(asset.Path)
		if err != nil {
			return fmt.Errorf("failed to open asset file %s: %w", asset.Path, err)
		}

		// Upload the asset
		err = r.UploadReleaseAsset(ctx, client, releaseID, asset.Name, mediaType, file)
		file.Close()
		if err != nil {
			return fmt.Errorf("failed to upload asset %s: %w", asset.Name, err)
		}
	}

	fmt.Printf("Successfully created release v%s with all assets\n", version)
	return nil
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

	// Create a context for GitHub operations
	ctx := context.Background()

	release := &ReleaseImpl{
		version:  version,
		repoPath: repoPath,
	}

	// Execute the release strategy
	if err := ReleaseStrategy(ctx, release, dirPath, pat); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute release strategy: %v\n", err)
		os.Exit(1)
	}
}

type HeaderRoundTripper struct {
	Transport http.RoundTripper
	Token     string
}

func (h *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	return h.Transport.RoundTrip(req)
}

func AuthenticatedHttpClient(token string) *http.Client {
	return &http.Client{
		Transport: &HeaderRoundTripper{
			Transport: http.DefaultTransport,
			Token:     token,
		},
	}
}
