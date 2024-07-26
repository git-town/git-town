package shared

import (
	"reflect"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

func OpcodeBranches(code Opcode) []gitdomain.BranchName {
	result := []gitdomain.BranchName{}

	valueOfSelf := reflect.ValueOf(code).Elem()
	typeOfSelf := valueOfSelf.Type()

	fieldCount := valueOfSelf.NumField()
	for i := 0; i < fieldCount; i++ {
		field := valueOfSelf.Field(i)
		fieldType := typeOfSelf.Field(i).Type

		// Check if the field is of one of the target types
		if fieldType == reflect.TypeOf((*gitdomain.BranchName)(nil)).Elem() {
			branchName := field.Interface().(gitdomain.BranchName)
			result = append(result, branchName)
		}
		if fieldType == reflect.TypeOf((*gitdomain.LocalBranchName)(nil)).Elem() {
			localBranchName := field.Interface().(gitdomain.LocalBranchName)
			result = append(result, localBranchName.BranchName())
		}
		if fieldType == reflect.TypeOf((*gitdomain.RemoteBranchName)(nil)).Elem() {
			remoteBranchName := field.Interface().(gitdomain.RemoteBranchName)
			result = append(result, remoteBranchName.BranchName())
		}
	}

	return result
}
