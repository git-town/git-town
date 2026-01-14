package mockproposals

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// MockProposalPath is the path to the mock proposals file.
type MockProposalPath string

func (self MockProposalPath) String() string {
	return string(self)
}

func Load(path MockProposalPath) MockProposals {
	fmt.Println("222222222222222222222222222222222222222222222222222", path)
	fileData, hasMockProposals := LoadBytes(path).Get()
	if !hasMockProposals {
		fmt.Println("333333333333333333333333333333333333333333333333333")
		return []forgedomain.ProposalData{}
	}
	fmt.Println("4444444444444444444444444444444444444444444444444444444444444")
	var proposals []forgedomain.ProposalData
	asserts.NoError(json.Unmarshal(fileData, &proposals))
	return proposals
}

func LoadBytes(path MockProposalPath) Option[[]byte] {
	fileData, err := os.ReadFile(path.String())
	if os.IsNotExist(err) {
		return None[[]byte]()
	}
	if err != nil {
		panic(err)
	}
	return Some(fileData)
}

func NewMockProposalPath(configDir configdomain.RepoConfigDir) MockProposalPath {
	return MockProposalPath(filepath.Join(configDir.String(), "proposals.json"))
}

func Save(path MockProposalPath, proposals MockProposals) string {
	fmt.Println("11111111111111111111111111111111111111111111111", path)
	asserts.NoError(os.MkdirAll(filepath.Dir(path.String()), 0o755))
	content := asserts.NoError1(json.MarshalIndent(proposals, "", "  "))
	asserts.NoError(os.WriteFile(path.String(), content, 0o600))
	return string(content)
}
