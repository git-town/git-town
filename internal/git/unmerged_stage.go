package git

// UnmergedStage describes the roles that a file can play in a merge conflict.
type UnmergedStage int

const (
	UnmergedStageBase          UnmergedStage = 1 // the base version in a 3-way merge
	UnmergedStageCurrentBranch UnmergedStage = 2 // the file on the branch on which the merge conflict happens
	UnmergedStageIncoming      UnmergedStage = 3 // the file on the branch getting merged in
)

// UnmergedStages provides all possible UnmergedStage instances.
var UnmergedStages = []UnmergedStage{
	UnmergedStageBase,
	UnmergedStageCurrentBranch,
	UnmergedStageIncoming,
}
