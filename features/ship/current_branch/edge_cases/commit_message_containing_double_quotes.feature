Feature: commit message with double-quotes

  Background:
    Given a Git repo clone
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    When I run "git-town ship -m 'with "double quotes"'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                              |
      | feature | git fetch --prune --tags             |
      |         | git checkout main                    |
      | main    | git merge --squash --ff feature      |
      |         | git commit -m "with "double quotes"" |
      |         | git push                             |
      |         | git push origin :feature             |
      |         | git branch -D feature                |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no uncommitted files exist
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | with "double quotes" |
    And no lineage exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                       |
      | main   | git revert {{ sha 'with "double quotes"' }}                   |
      |        | git push                                                      |
      |        | git push origin {{ sha 'initial commit' }}:refs/heads/feature |
      |        | git branch feature {{ sha 'feature commit' }}                 |
      |        | git checkout feature                                          |
    And the current branch is now "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE                       |
      | main    | local, origin | with "double quotes"          |
      |         |               | Revert "with "double quotes"" |
      | feature | local         | feature commit                |
    And the initial branches and lineage exist
