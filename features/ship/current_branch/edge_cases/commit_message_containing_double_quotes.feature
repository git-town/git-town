Feature: commit message with double-quotes

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    When I run "git-town ship -m 'with "double quotes"'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                              |
      | feature | git fetch --prune --tags             |
      |         | git checkout main                    |
      | main    | git merge --squash feature           |
      |         | git commit -m "with "double quotes"" |
      |         | git push                             |
      |         | git push origin :feature             |
      |         | git branch -D feature                |
    And the current branch is now "main"
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And no uncommitted files exist
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | with "double quotes" |
    And no branch hierarchy exists now

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
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                       |
      | main    | local, origin | with "double quotes"          |
      |         |               | Revert "with "double quotes"" |
      | feature | local         | feature commit                |
    And the initial branches and hierarchy exist
