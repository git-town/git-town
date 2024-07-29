package shared

import (
	"reflect"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

func BranchesInOpcode(code Opcode) []gitdomain.BranchName {
	result := []gitdomain.BranchName{}
	valueOfSelf := reflect.ValueOf(code).Elem()
	typeOfSelf := valueOfSelf.Type()
	for i := range valueOfSelf.NumField() {
		field := valueOfSelf.Field(i)
		fieldType := typeOfSelf.Field(i).Type
		if fieldType == reflect.TypeOf((*gitdomain.BranchName)(nil)).Elem() {
			branchName := field.Interface().(gitdomain.BranchName)
			result = append(result, branchName)
		}
		if fieldType == reflect.TypeOf((*gitdomain.LocalBranchName)(nil)).Elem() {
			localBranchName := field.Interface().(gitdomain.LocalBranchName)
			result = append(result, localBranchName.BranchName())
		}
		if fieldType == reflect.TypeOf((*gitdomain.LocalBranchNames)(nil)).Elem() {
			localBranchNames := field.Interface().(gitdomain.LocalBranchNames)
			result = append(result, localBranchNames.BranchNames()...)
		}
		if fieldType == reflect.TypeOf((*gitdomain.RemoteBranchName)(nil)).Elem() {
			remoteBranchName := field.Interface().(gitdomain.RemoteBranchName)
			result = append(result, remoteBranchName.BranchName())
		}
	}
	return result
}
