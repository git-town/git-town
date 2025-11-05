package shared

import (
	"reflect"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

func BranchesInOpcode(code Opcode) []gitdomain.BranchName {
	var result []gitdomain.BranchName
	valueOfSelf := reflect.ValueOf(code).Elem()
	typeOfSelf := valueOfSelf.Type()
	for i := range valueOfSelf.NumField() {
		field := valueOfSelf.Field(i)
		fieldType := typeOfSelf.Field(i).Type
		if fieldType == reflect.TypeFor[gitdomain.BranchName]() {
			branchName := field.Interface().(gitdomain.BranchName)
			result = append(result, branchName)
		}
		if fieldType == reflect.TypeFor[gitdomain.LocalBranchName]() {
			localBranchName := field.Interface().(gitdomain.LocalBranchName)
			result = append(result, localBranchName.BranchName())
		}
		if fieldType == reflect.TypeFor[gitdomain.LocalBranchNames]() {
			localBranchNames := field.Interface().(gitdomain.LocalBranchNames)
			result = append(result, localBranchNames.BranchNames()...)
		}
		if fieldType == reflect.TypeFor[gitdomain.RemoteBranchName]() {
			remoteBranchName := field.Interface().(gitdomain.RemoteBranchName)
			result = append(result, remoteBranchName.BranchName())
		}
	}
	return result
}
