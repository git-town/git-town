package mockproposals

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/pkg/asserts"
)

// MockProposalPath is the path to the mock proposals file.
type MockProposalPath string

func NewMockProposalPath(configDir configdomain.RepoConfigDir) MockProposalPath {
	return MockProposalPath(filepath.Join(configDir.String(), "proposals.json"))
}

func Load(path MockProposalPath) MockProposals {
	fileData := LoadBytes(workspaceDir)
	var proposals []forgedomain.ProposalData
	asserts.NoError(json.Unmarshal(fileData, &proposals))
	return MockProposals{
		Dir:       workspaceDir,
		Proposals: proposals,
	}
}

func LoadBytes(workspaceDir string) []byte {
	filePath := FilePath(workspaceDir)
	fileData := asserts.NoError1(os.ReadFile(filePath))
	return fileData
}

func Save(workspaceDir string, proposals []forgedomain.ProposalData) string {
	content := asserts.NoError1(json.MarshalIndent(proposals, "", "  "))
	asserts.NoError(os.WriteFile(FilePath(workspaceDir), content, 0o600))
	return string(content)
}
