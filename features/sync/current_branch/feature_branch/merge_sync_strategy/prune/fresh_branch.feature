Feature: prune a freshly created branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    When I run "git-town sync --prune"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git push origin :feature |
      |         | git branch -D feature    |
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                         |
      | feature-2 | git push origin {{ sha 'initial commit' }}:refs/heads/feature-1 |
      |           | git branch feature-1 {{ sha 'feature-1 commit' }}               |
      |           | git checkout feature-1                                          |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, origin | main commit      |
      | feature-1 | local         | feature-1 commit |
      | feature-2 | local, origin | feature-2 commit |
