package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-github/v69/github"
	"go.uber.org/mock/gomock"
)

// TestReleaseStrategy_HappyPath verifies the correct sequence of operations
// in the release strategy. Update this test when adding new strategy steps.
func TestReleaseStrategy_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockRelease(ctrl)
	testDir := "/test/path"
	testRepo := "/test/repo"
	testPAT := "test_pat"
	testVersion := "1.2.3"
	testReleaseID := int64(123456)
	testNotes := "Test release notes"
	testFormula := "test formula content"
	testCommitSHA := "abcdef1234567890"
	testAssets := []ReleaseAsset{
		{Name: "asset1", Path: "/test/path/asset1.zip", Checksum: "checksum1"},
		{Name: "asset2", Path: "/test/path/asset2.zip", Checksum: "checksum2"},
	}

	// Mock the GitHub client
	mockClient := &github.Client{}

	// Create a context for the test
	ctx := context.Background()

	// Set up mock expectations for GetRepoPath and GetVersion that can be called multiple times
	mock.EXPECT().GetRepoPath().Return(testRepo).AnyTimes()
	mock.EXPECT().GetVersion().Return(testVersion).AnyTimes()
	mock.EXPECT().GetAssets().Return(testAssets).AnyTimes()
	mock.EXPECT().GetNotesContent().Return(testNotes).AnyTimes()

	// Create a temporary file for asset mocking
	tmpFile1, err := os.CreateTemp("", "test-asset1")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile1.Name())
	defer tmpFile1.Close()

	tmpFile2, err := os.CreateTemp("", "test-asset2")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile2.Name())
	defer tmpFile2.Close()

	// Write some dummy content
	if _, err := tmpFile1.WriteString("test content 1"); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if _, err := tmpFile2.WriteString("test content 2"); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Reset file position to beginning
	if _, err := tmpFile1.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}
	if _, err := tmpFile2.Seek(0, 0); err != nil {
		t.Fatalf("Failed to seek temp file: %v", err)
	}

	// Expected sequence of operations
	gomock.InOrder(
		// Step 1: Create GitHub client
		mock.EXPECT().CreateGitHubClient(testPAT).Return(mockClient),

		// Step 2: Download assets and release notes
		mock.EXPECT().DownloadReleaseAssetsAndNotes(ctx, mockClient, testDir, testPAT).Return(nil),

		// Step 3: Validate assets
		mock.EXPECT().Validate().Return(nil),
		mock.EXPECT().ComputeChecksums().Return(nil),

		// Step 4: Check if release exists
		mock.EXPECT().ReleaseExists(ctx, mockClient, testVersion).Return(false, nil),

		// Step 5: Generate and save the Homebrew formula
		mock.EXPECT().RenderFormulaTemplate().Return(testFormula, nil),
		mock.EXPECT().SaveFormulaFile(testFormula).Return(nil),

		// Step 6: Commit the formula change
		mock.EXPECT().CommitFormulaChange(ctx, testPAT, gomock.Any()).Return(testCommitSHA, nil),

		// Step 7: Create and push tag
		mock.EXPECT().CreateAndPushTagToGitHub(ctx, testPAT, testVersion, testRepo, testCommitSHA).Return(nil),

		// Step 8: Create GitHub release
		mock.EXPECT().CreateGitHubRelease(ctx, mockClient, testVersion, testNotes).Return(testReleaseID, nil),

		// Step 9: Upload assets (one per asset)
		mock.EXPECT().OpenAssetFile(testAssets[0].Path).Return(tmpFile1, nil),
		mock.EXPECT().UploadReleaseAsset(ctx, mockClient, testReleaseID, testAssets[0].Name, mediaType, tmpFile1).Return(nil),

		mock.EXPECT().OpenAssetFile(testAssets[1].Path).Return(tmpFile2, nil),
		mock.EXPECT().UploadReleaseAsset(ctx, mockClient, testReleaseID, testAssets[1].Name, mediaType, tmpFile2).Return(nil),

		// Step 10: Deploy to Vercel
		mock.EXPECT().DeployToVercel(ctx).Return(nil),
	)

	// Execute the strategy
	if err := ReleaseStrategy(ctx, mock, testDir, testPAT); err != nil {
		t.Errorf("ReleaseStrategy failed: %v", err)
	}
}

// TestReleaseStrategy_ReleaseAlreadyExists tests that the strategy fails correctly when a release already exists.
func TestReleaseStrategy_ReleaseAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockRelease(ctrl)
	testDir := "/test/path"
	testRepo := "/test/repo"
	testPAT := "test_pat"
	testVersion := "1.2.3"
	testAssets := []ReleaseAsset{
		{Name: "asset1", Path: "/test/path/asset1.zip", Checksum: "checksum1"},
	}

	// Mock the GitHub client
	mockClient := &github.Client{}

	// Create a context for the test
	ctx := context.Background()

	// Set up mock expectations for GetRepoPath and GetVersion that can be called multiple times
	mock.EXPECT().GetRepoPath().Return(testRepo).AnyTimes()
	mock.EXPECT().GetVersion().Return(testVersion).AnyTimes()
	mock.EXPECT().GetAssets().Return(testAssets).AnyTimes()

	// Expected sequence of operations
	gomock.InOrder(
		// Step 1: Create GitHub client
		mock.EXPECT().CreateGitHubClient(testPAT).Return(mockClient),

		// Step 2: Download assets and release notes
		mock.EXPECT().DownloadReleaseAssetsAndNotes(ctx, mockClient, testDir, testPAT).Return(nil),

		// Step 3: Validate assets
		mock.EXPECT().Validate().Return(nil),
		mock.EXPECT().ComputeChecksums().Return(nil),

		// Step 4: Check if release exists - this time it DOES exist
		mock.EXPECT().ReleaseExists(ctx, mockClient, testVersion).Return(true, nil),
	)

	// Execute the strategy
	err := ReleaseStrategy(ctx, mock, testDir, testPAT)

	// Verify that we get an error about the release already existing
	if err == nil {
		t.Error("ReleaseStrategy should fail when release already exists")
	}

	expectedErrMsg := fmt.Sprintf("release %s already exists", tagName(testVersion))
	if err != nil && err.Error() != expectedErrMsg {
		t.Errorf("Expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}

func TestReleaseImpl_GetVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "normal version",
			version: "1.2.3",
			want:    "1.2.3",
		},
		{
			name:    "empty version",
			version: "",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReleaseImpl{
				version: tt.version,
			}
			if got := r.GetVersion(); got != tt.want {
				t.Errorf("ReleaseImpl.GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReleaseImpl_GetAssets(t *testing.T) {
	assets := []ReleaseAsset{
		{Name: "test1", Path: "/path/1", Checksum: "abc"},
		{Name: "test2", Path: "/path/2", Checksum: "def"},
	}

	r := &ReleaseImpl{
		assets: assets,
	}

	got := r.GetAssets()
	if len(got) != len(assets) {
		t.Errorf("ReleaseImpl.GetAssets() returned %d assets, want %d", len(got), len(assets))
	}

	for i, asset := range got {
		if asset != assets[i] {
			t.Errorf("ReleaseImpl.GetAssets()[%d] = %v, want %v", i, asset, assets[i])
		}
	}
}

func TestReleaseImpl_GetNotesContent(t *testing.T) {
	content := "Test release notes\nWith multiple lines"
	r := &ReleaseImpl{
		notesContent: content,
	}

	if got := r.GetNotesContent(); got != content {
		t.Errorf("ReleaseImpl.GetNotesContent() = %v, want %v", got, content)
	}
}

func TestReleaseImpl_GetRepoPath(t *testing.T) {
	path := "/test/repo/path"
	r := &ReleaseImpl{
		repoPath: path,
	}

	if got := r.GetRepoPath(); got != path {
		t.Errorf("ReleaseImpl.GetRepoPath() = %v, want %v", got, path)
	}
}

func TestReleaseImpl_OpenAssetFile(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test-asset-file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write test content
	testContent := "test file content"
	if _, err := tmpFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close() // Close the file before reopening

	// Use the OpenAssetFile method to open the file
	r := &ReleaseImpl{}
	file, err := r.OpenAssetFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReleaseImpl.OpenAssetFile() error = %v", err)
	}
	defer file.Close()

	// Read the content to verify it works
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Failed to read from opened file: %v", err)
	}

	// Verify content matches
	if string(content) != testContent {
		t.Errorf("File content = %q, want %q", string(content), testContent)
	}
}

func TestTagName(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "normal version",
			version: "1.2.3",
			want:    "v1.2.3",
		},
		{
			name:    "already has v prefix",
			version: "v1.2.3", // Should still work even though the 'v' will be stripped elsewhere
			want:    "vv1.2.3",
		},
		{
			name:    "empty version",
			version: "",
			want:    "v",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tagName(tt.version); got != tt.want {
				t.Errorf("tagName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReleaseImpl_CreateAssets(t *testing.T) {
	r := &ReleaseImpl{}
	testVersion := "1.2.3"
	err := r.CreateAssets("/test/path", testVersion)
	if err != nil {
		t.Errorf("ReleaseImpl.CreateAssets() unexpected error = %v", err)
	}

	assets := r.GetAssets()
	expectedCount := 4 // We have 4 platform combinations
	if len(assets) != expectedCount {
		t.Errorf("ReleaseImpl.CreateAssets() created %d assets, want %d", len(assets), expectedCount)
	}

	// Check that all expected asset names are present with version
	expectedNames := map[string]bool{
		assetNameMacosArm64.String(testVersion): false,
		assetNameMacosIntel.String(testVersion): false,
		assetNameLinuxArm64.String(testVersion): false,
		assetNameLinuxIntel.String(testVersion): false,
	}

	for _, asset := range assets {
		if _, ok := expectedNames[asset.Name]; !ok {
			t.Errorf("Unexpected asset name: %s", asset.Name)
			continue
		}
		expectedNames[asset.Name] = true
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Missing expected asset: %s", name)
		}
	}

	// Verify the paths are constructed correctly
	for _, asset := range assets {
		expectedPath := "/test/path/" + asset.Name + "." + fileType
		if asset.Path != expectedPath {
			t.Errorf("Asset %s has path %s, want %s", asset.Name, asset.Path, expectedPath)
		}
	}
}

func TestReleaseImpl_RenderFormulaTemplate(t *testing.T) {
	// Create a ReleaseImpl with assets and checksums
	version := "1.2.3"
	r := &ReleaseImpl{
		version: version,
		assets: []ReleaseAsset{
			{Name: assetNameMacosArm64.String(version), Path: "/path/darwin-arm64.zip", Checksum: "mac-arm-checksum"},
			{Name: assetNameMacosIntel.String(version), Path: "/path/darwin-amd64.zip", Checksum: "mac-intel-checksum"},
			{Name: assetNameLinuxArm64.String(version), Path: "/path/linux-arm64.zip", Checksum: "linux-arm-checksum"},
			{Name: assetNameLinuxIntel.String(version), Path: "/path/linux-amd64.zip", Checksum: "linux-intel-checksum"},
		},
	}

	// Render the template
	output, err := r.RenderFormulaTemplate()
	if err != nil {
		t.Fatalf("RenderFormulaTemplate failed: %v", err)
	}

	// Check that the output contains the expected values
	expectedValues := []string{
		"version \"1.2.3\"",
		"sha256 \"mac-arm-checksum\"",
		"sha256 \"mac-intel-checksum\"",
		"sha256 \"linux-arm-checksum\"",
		"sha256 \"linux-intel-checksum\"",
	}

	for _, expectedValue := range expectedValues {
		if !strings.Contains(output, expectedValue) {
			t.Errorf("Expected formula to contain %q, but it didn't", expectedValue)
		}
	}
}

func TestReleaseImpl_GetFormulaFilePath(t *testing.T) {
	testPath := "/test/repo/path"
	r := &ReleaseImpl{
		repoPath: testPath,
	}

	expected := filepath.Join(testPath, "Formula", formulaFileName)
	actual := r.GetFormulaFilePath()

	if actual != expected {
		t.Errorf("GetFormulaFilePath() = %v, want %v", actual, expected)
	}
}

// TestReleaseImpl_DownloadReleaseAssetsAndNotes tests the downloading of release assets from GitHub
func TestReleaseImpl_DownloadReleaseAssetsAndNotes(t *testing.T) {
	// We'll skip this test as it requires more complex setup that we'll handle later
	t.Skip("This test requires more complex setup and will be implemented later")
}

// TestDeployToVercel tests the Vercel deployment functionality
func TestDeployToVercel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Create a temporary directory to simulate a repository
	tmpRepo, err := os.MkdirTemp("", "test-repo")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpRepo)

	// Create a release implementation with the test repo path
	r := &ReleaseImpl{
		version:  "1.2.3",
		repoPath: tmpRepo,
	}

	// Create a mock vercel command for testing
	// In a real test, you would use a mock command runner or dependency injection
	// Here we're just testing the execution logic by creating a mock script

	mockVercelPath := filepath.Join(tmpRepo, "mock-vercel")
	mockScript := `#!/bin/sh
echo "Mock Vercel executed with args: $@"
exit 0
`
	if err := os.WriteFile(mockVercelPath, []byte(mockScript), 0755); err != nil {
		t.Fatalf("Failed to write mock vercel script: %v", err)
	}

	// Back up the original exec.Command function and restore it after the test
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	// Mock the exec.Command to use our script instead of the real vercel
	execCommand = func(name string, args ...string) *exec.Cmd {
		if strings.Contains(name, "vercel") {
			return exec.Command(mockVercelPath, args...)
		}
		return exec.Command(name, args...)
	}

	// Execute the deployment
	ctx := context.Background()
	err = r.DeployToVercel(ctx)

	// Check the result
	if err != nil {
		t.Errorf("DeployToVercel failed: %v", err)
	}
}
