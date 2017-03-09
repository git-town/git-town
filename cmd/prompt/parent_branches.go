package prompt

import (
  "fmt"

  "github.com/Originate/gt/cmd/config"
  "github.com/Originate/gt/cmd/util"

  "github.com/fatih/color"
)

func EnsureKnowsParentBranches(branchNames []string) {
  headerShown := false
  for _, branchName := range(branchNames) {
    if config.KnowsAllAncestorBranches(branchName) {
      continue
    }
    if !headerShown {
      printParentBranchHeader()
      headerShown = true
    }
    askForBranchAncestry(branchName)
    ancestors := config.CompileAncestorBranches(branchName)
    config.SetAncestorBranches(branchName, ancestors)
  }
  if headerShown {
    fmt.Println()
  }
}

func askForBranchAncestry(branchName string) {
  current := branchName
  for {
    fmt.Println("parent", config.GetParentBranch(current))
    if config.GetParentBranch(current) != "" {
      break
    }
    parent := askForParentBranch(current)
    config.SetParentBranch(current, parent)
    if parent == config.GetMainBranch() {
      break
    }
    current = parent
  }
}

func askForParentBranch(branchName string) string {
  fmt.Printf(
    "Please specify the parent branch of %s by name or number (default: %s): ",
    color.New(color.FgCyan).Sprintf(branchName),
    config.GetMainBranch(),
  )
  return util.ReadLineFromStdin()
}

func printParentBranchHeader() {

}

// # Makes sure that we know all the parent branches
// # Asks the user if necessary
// function ensure_knows_parent_branches {
//   local branches=$1 # space separated list of branches
//
//   local ancestors
//   local branch
//   local child
//   local header_shown=false
//   local numerical_regex='^[0-9]+$'
//   local parent
//   local user_input
//
//   for branch in $branches; do
//     child=$branch
//     if [ "$(knows_all_ancestor_branches "$child")" = true ]; then
//       continue
//     fi
//     if [ "$(is_perennial_branch "$child")" = true ]; then
//       continue
//     fi
//
//     while [ "$child" != "$MAIN_BRANCH_NAME" ]; do
//       if [ "$(knows_parent_branch "$child")" = true ]; then
//         parent=$(parent_branch "$child")
//       else
//         if [ "$header_shown" = false ]; then
//           echo_parent_branch_header
//           header_shown=true
//         fi
//
//         parent=""
//         while [ -z "$parent" ]; do
//           echo -n "Please specify the parent branch of $(echo_inline_cyan_bold "$child") by name or number (default: $MAIN_BRANCH_NAME): "
//           read user_input
//           if [[ $user_input =~ $numerical_regex ]] ; then
//             # user entered a number here
//             parent="$(get_numbered_branch "$user_input")"
//             if [ -z "$parent" ]; then
//               echo_error_header
//               echo_error "Invalid branch number"
//             fi
//           elif [ -z "$user_input" ]; then
//             # user entered nothing
//             parent=$MAIN_BRANCH_NAME
//           else
//             if [ "$(has_branch "$user_input")" == true ]; then
//               parent=$user_input
//             else
//               echo_error_header
//               echo_error "Branch '$user_input' doesn't exist"
//             fi
//           fi
//           if [ "$child" = "$parent" ]; then
//             echo_error_header
//             echo_error "'$child' cannot be the parent of itself"
//             parent=''
//           elif [ "$(has_ancestor_branch "$parent" "$child")" = true ]; then
//             echo_error_header
//             echo_error "Nested branch loop detected: '$child' is an ancestor of '$parent'"
//             parent=''
//           fi
//         done
//         store_parent_branch "$child" "$parent"
//       fi
//       child=$parent
//     done
//     ancestors=$(compile_ancestor_branches "$branch")
//     store_ancestor_branches "$branch" "$ancestors"
//   done
//
//   if [ "$header_shown" = true ]; then
//     echo
//   fi
// }
