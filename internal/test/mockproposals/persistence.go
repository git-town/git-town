package mockproposals

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/pkg/asserts"
)

func FilePath(workspaceDir string) string {
	return filepath.Join(workspaceDir, "proposals.json")
}

func Load(workspaceDir string) MockProposals {
	proposals, err := os.ReadFile(FilePath(workspaceDir))
	if err != nil {
		return MockProposals{}
	}
	var result MockProposals
	asserts.NoError(json.Unmarshal(proposals, &result))
	return result
}

func Save(workspaceDir string, proposals MockProposals) {
	content := asserts.NoError1(json.MarshalIndent(proposals, "", "  "))
	fmt.Println("22222222222222222222222222222", string(content))
	asserts.NoError(os.WriteFile(FilePath(workspaceDir), content, 0o600))
}
