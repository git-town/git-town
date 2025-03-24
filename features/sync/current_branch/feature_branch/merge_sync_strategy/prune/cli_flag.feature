Feature: prune enabled via CLI

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME | FILE CONTENT |
      | main    | local    | main commit    | file      | content      |
      | feature | local    | feature commit | file      | content      |
    And the current branch is "feature"
    When I run "git-town sync --prune"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main --no-update-refs |
      |         | git push                                |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff main           |
      |         | git merge --no-edit --ff origin/feature |
      |         | git checkout main                       |
      | main    | git push origin :feature                |
      |         | git branch -D feature                   |
    And all branches are now synchronized
    And the current branch is now "main"
    And these branches exist now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                       |
      | main   | git push origin {{ sha 'initial commit' }}:refs/heads/feature |
      |        | git branch feature {{ sha 'feature commit' }}                 |
      |        | git checkout feature                                          |
    And the current branch is now "feature"
    And these branches exist now
      | REPOSITORY    | BRANCHES      |
      | local, origin | main, feature |
