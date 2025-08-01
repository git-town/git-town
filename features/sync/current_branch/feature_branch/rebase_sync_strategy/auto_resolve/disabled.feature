Feature: don't auto-resolve merge conflicts

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE    | PARENT | LOCATIONS     |
      | child      | feature | main   | local, origin |
      | grandchild | feature | child  | local, origin |
    And the commits
      | BRANCH     | LOCATION | MESSAGE                       | FILE NAME        | FILE CONTENT       |
      | main       | local    | conflicting main commit       | conflicting_file | main content       |
      | child      | local    | child commit                  | child_file       | child content      |
      | grandchild | local    | conflicting grandchild commit | conflicting_file | grandchild content |
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And origin deletes the "child" branch
    And the current branch is "child" and the previous branch is "grandchild"
    When I run "git-town sync --all --auto-resolve=0"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH     | COMMAND                                                 |
      | child      | git fetch --prune --tags                                |
      |            | git checkout main                                       |
      | main       | git -c rebase.updateRefs=false rebase origin/main       |
      |            | git push                                                |
      |            | git checkout grandchild                                 |
      | grandchild | git pull                                                |
      |            | git -c rebase.updateRefs=false rebase --onto main child |
      |            | GIT_EDITOR=true git rebase --continue                   |
      |            | git push --force-with-lease                             |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And Git Town prints something like:
      """
      could not apply .* conflicting grandchild commit
      """
    And a rebase is now in progress
