package main

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestUncachedFilePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		cachedFile string
		pkgPath    string
		want       string
	}{
		{
			name:       "cached_api_connector.go",
			cachedFile: "cached_api_connector.go",
			pkgPath:    "internal/forge/gitlab",
			want:       "internal/forge/gitlab/api_connector.go",
		},
		{
			name:       "cached_connector.go",
			cachedFile: "cached_connector.go",
			pkgPath:    "internal/forge/glab",
			want:       "internal/forge/glab/connector.go",
		},
		{
			name:       "full path with cached_api_connector.go",
			cachedFile: "/home/user/project/internal/forge/gitlab/cached_api_connector.go",
			pkgPath:    "internal/forge/gitlab",
			want:       "internal/forge/gitlab/api_connector.go",
		},
		{
			name:       "file without cached_ prefix",
			cachedFile: "api_connector.go",
			pkgPath:    "internal/forge/gitlab",
			want:       "internal/forge/gitlab/api_connector.go",
		},
		{
			name:       "empty package path",
			cachedFile: "cached_connector.go",
			pkgPath:    "",
			want:       "connector.go",
		},
		{
			name:       "multiple cached_ prefixes",
			cachedFile: "cached_cached_connector.go",
			pkgPath:    "internal/forge/gitlab",
			want:       "internal/forge/gitlab/cached_connector.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			have := uncachedFilePath(tt.cachedFile, tt.pkgPath)
			must.EqOp(t, tt.want, have)
		})
	}
}
