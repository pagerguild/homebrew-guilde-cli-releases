package main

import (
	"testing"

	"go.uber.org/mock/gomock"
)

// TestReleaseStrategy_HappyPath verifies the correct sequence of operations
// in the release strategy. Update this test when adding new strategy steps.
func TestReleaseStrategy_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := NewMockRelease(ctrl)
	testDir := "/test/path"

	// Expected sequence of operations
	gomock.InOrder(
		mock.EXPECT().CreateAssets(testDir).Return(nil),
		mock.EXPECT().Validate().Return(nil),
		mock.EXPECT().ComputeChecksums().Return(nil),
		mock.EXPECT().LoadReleaseNotes().Return(nil),
	)

	if err := ReleaseStrategy(mock, testDir); err != nil {
		t.Errorf("ReleaseStrategy failed: %v", err)
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

func TestReleaseImpl_CreateAssets(t *testing.T) {
	r := &ReleaseImpl{}
	err := r.CreateAssets("/test/path")
	if err != nil {
		t.Errorf("ReleaseImpl.CreateAssets() unexpected error = %v", err)
	}

	assets := r.GetAssets()
	expectedCount := 4 // We have 4 platform combinations
	if len(assets) != expectedCount {
		t.Errorf("ReleaseImpl.CreateAssets() created %d assets, want %d", len(assets), expectedCount)
	}

	// Check that all expected asset names are present
	expectedNames := map[string]bool{
		assetNameMacosArm64: false,
		assetNameMacosIntel: false,
		assetNameLinuxArm64: false,
		assetNameLinuxIntel: false,
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
