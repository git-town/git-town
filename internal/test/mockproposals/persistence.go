package mockproposals

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/pkg/asserts"
)

func FilePath(workspaceDir string) string {
	return filepath.Join(workspaceDir, "proposals.json")
}

func Load(workspaceDir string) MockProposals {
	filePath := FilePath(workspaceDir)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return MockProposals{
			Dir:       workspaceDir,
			Proposals: []forgedomain.ProposalData{},
		}
	}
	var proposals []forgedomain.ProposalData
	asserts.NoError(json.Unmarshal(fileData, &proposals))
	return MockProposals{
		Dir:       workspaceDir,
		Proposals: proposals,
	}
}

func Save(workspaceDir string, proposals []forgedomain.ProposalData) {
	content := asserts.NoError1(json.MarshalIndent(proposals, "", "  "))
	asserts.NoError(os.WriteFile(FilePath(workspaceDir), content, 0o600))
}
