package format

import "testing"

func TestBuildCollectionPath(t *testing.T) {
	tests := []struct {
		name         string
		collectionID string
		itemID       string
		extension    string
		want         string
	}{
		{
			name:         "markdown rule",
			collectionID: "cleanCode",
			itemID:       "meaningfulNames",
			extension:    ".md",
			want:         "cleanCode_meaningfulNames.md",
		},
		{
			name:         "copilot prompt",
			collectionID: "codeReview",
			itemID:       "reviewPR",
			extension:    ".prompt.md",
			want:         "codeReview_reviewPR.prompt.md",
		},
		{
			name:         "cursor rule",
			collectionID: "cleanCode",
			itemID:       "smallFunctions",
			extension:    ".mdc",
			want:         "cleanCode_smallFunctions.mdc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildCollectionPath(tt.collectionID, tt.itemID, tt.extension)
			if got != tt.want {
				t.Errorf("BuildCollectionPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildStandalonePath(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		extension  string
		want       string
	}{
		{
			name:       "markdown rule",
			resourceID: "meaningfulNames",
			extension:  ".md",
			want:       "meaningfulNames.md",
		},
		{
			name:       "copilot prompt",
			resourceID: "reviewPR",
			extension:  ".prompt.md",
			want:       "reviewPR.prompt.md",
		},
		{
			name:       "cursor rule",
			resourceID: "smallFunctions",
			extension:  ".mdc",
			want:       "smallFunctions.mdc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildStandalonePath(tt.resourceID, tt.extension)
			if got != tt.want {
				t.Errorf("BuildStandalonePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildClaudeCollectionPath(t *testing.T) {
	tests := []struct {
		name         string
		collectionID string
		itemID       string
		want         string
	}{
		{
			name:         "claude prompt from spec",
			collectionID: "codeReview",
			itemID:       "reviewPR",
			want:         "codeReview_reviewPR/SKILL.md",
		},
		{
			name:         "another collection",
			collectionID: "testing",
			itemID:       "unitTests",
			want:         "testing_unitTests/SKILL.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildClaudeCollectionPath(tt.collectionID, tt.itemID)
			if got != tt.want {
				t.Errorf("BuildClaudeCollectionPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildClaudeStandalonePath(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		want       string
	}{
		{
			name:       "claude prompt from spec",
			resourceID: "reviewPR",
			want:       "reviewPR/SKILL.md",
		},
		{
			name:       "another standalone",
			resourceID: "unitTests",
			want:       "unitTests/SKILL.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildClaudeStandalonePath(tt.resourceID)
			if got != tt.want {
				t.Errorf("BuildClaudeStandalonePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
