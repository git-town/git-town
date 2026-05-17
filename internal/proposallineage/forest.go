package proposallineage

type Forest []TreeNode

func (self Forest) BranchCount() int {
	var count int
	for _, node := range self {
		count += node.BranchCount()
	}
	return count
}
