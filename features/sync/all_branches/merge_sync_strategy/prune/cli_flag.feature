Feature: prune enabled via CLI

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE    | PARENT    | LOCATIONS     |
      | feature-1 | feature | main      | local, origin |
      | feature-2 | feature | feature-1 | local, origin |
      | feature-3 | feature | feature-2 | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          | FILE NAME  | FILE CONTENT  |
      | main      | local         | main commit      | file       | content       |
      | feature-1 | local         | feature-1 commit | file       | content       |
      | feature-2 | origin        | feature-2 commit | file       | content       |
      | feature-3 | local, origin | feature-3 commit | other_file | other content |
    And the current branch is "feature-1"
    When I run "git-town sync --all --prune"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND |
      | feature-1 | git fetch --prune --tags |
      | | git checkout main |
      | main | git -c rebase.updateRefs=false rebase origin/main |
      | | git push |
      | | git checkout feature-1 |
      | feature-1 | git merge --no-edit --ff main |
      | | git merge --no-edit --ff origin/feature-1 |
      | | git checkout main |
      | main | git push origin :feature-1 |
      | | git branch -D feature-1 |
      | | git checkout feature-2 |
      | feature-2 | git merge --no-edit --ff main |
      | | git merge --no-edit --ff origin/feature-2 |
      | | git checkout main |
      | main | git push origin :feature-2 |
      | | git branch -D feature-2 |
      | | git checkout feature-3 |
      | feature-3 | git merge --no-edit --ff main |
      | | git push |
      | | git push --tags |
    And this lineage exists now
      """
      main
        feature-3
      """
    And the branches are now
      | REPOSITORY    | BRANCHES        |
      | local, origin | main, feature-3 |
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                                                     |
      | feature-3 | git reset --hard {{ sha 'feature-3 commit' }}                               |
      |           | git push --force-with-lease --force-if-includes                             |
      |           | git push origin {{ sha 'initial commit' }}:refs/heads/feature-1             |
      |           | git push origin {{ sha-in-origin 'feature-2 commit' }}:refs/heads/feature-2 |
      |           | git branch feature-1 {{ sha 'feature-1 commit' }}                           |
      |           | git branch feature-2 {{ sha 'initial commit' }}                             |
      |           | git checkout feature-1                                                      |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE          |
      | main      | local, origin | main commit      |
      | feature-1 | local         | feature-1 commit |
      | feature-2 | origin        | feature-2 commit |
      | feature-3 | local, origin | feature-3 commit |
