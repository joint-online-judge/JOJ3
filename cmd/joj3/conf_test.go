package main

import (
	"reflect"
	"testing"
)

func TestParseConventionalCommit(t *testing.T) {
	tests := []struct {
		name    string
		commit  string
		want    *ConventionalCommit
		wantErr bool
	}{
		{
			name:   "Simple feat commit",
			commit: "feat: add new feature",
			want: &ConventionalCommit{
				Type:        "feat",
				Description: "add new feature",
			},
			wantErr: false,
		},
		{
			name:   "Commit with scope",
			commit: "fix(core): resolve memory leak",
			want: &ConventionalCommit{
				Type:        "fix",
				Scope:       "core",
				Description: "resolve memory leak",
			},
			wantErr: false,
		},
		{
			name:   "Breaking change commit",
			commit: "feat(api)!: redesign user authentication",
			want: &ConventionalCommit{
				Type:        "feat",
				Scope:       "api",
				Description: "redesign user authentication",
			},
			wantErr: false,
		},
		{
			name:   "Commit with body",
			commit: "docs: update README\n\nAdd installation instructions and improve examples",
			want: &ConventionalCommit{
				Type:        "docs",
				Description: "update README",
				Body:        "Add installation instructions and improve examples",
			},
			wantErr: false,
		},
		{
			name:   "Full commit with body and footer",
			commit: "feat(auth)!: implement OAuth2\n\nThis commit adds OAuth2 support to the authentication system.\n\nBREAKING CHANGE: Previous authentication tokens are no longer valid.",
			want: &ConventionalCommit{
				Type:        "feat",
				Scope:       "auth",
				Description: "implement OAuth2",
				Body:        "This commit adds OAuth2 support to the authentication system.",
				Footer:      "BREAKING CHANGE: Previous authentication tokens are no longer valid.",
			},
			wantErr: false,
		},
		{
			name:    "Invalid commit format",
			commit:  "This is not a valid conventional commit",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConventionalCommit(tt.commit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConventionalCommit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConventionalCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}
