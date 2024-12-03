Feature: the branch was shipped manually on the local machine

  Background:
    Given a Git repo with origin
    And Git Town setting "default-branch-type" is "prototype"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "main"
    And I ran "git merge feature --squash"
    And I ran "git commit -m merged"
    And origin deletes the "feature" branch
    When I run "git-town sync --all"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                 |
      | main   | git fetch --prune --tags                |
      |        | git rebase origin/main --no-update-refs |
      |        | git push                                |
      |        | git branch -D feature                   |
      |        | git push --tags                         |
    And the current branch is still "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                       |
      | main   | git branch feature {{ sha 'feature commit' }} |
    And the current branch is still "main"
    And the branches are now
      | REPOSITORY | BRANCHES      |
      | local      | main, feature |
      | origin     | main          |
    And this lineage exists now
      | BRANCH  | PARENT |
      | feature | main   |
