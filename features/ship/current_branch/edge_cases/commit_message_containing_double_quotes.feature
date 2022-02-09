Feature: commit message can contain double-quotes

  Background:
    Given a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And I am on the "feature" branch
    When I run "git-town ship -m 'with "double quotes"'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                              |
      | feature | git fetch --prune --tags             |
      |         | git checkout main                    |
      | main    | git rebase origin/main               |
      |         | git checkout feature                 |
      | feature | git merge --no-edit origin/feature   |
      |         | git merge --no-edit main             |
      |         | git checkout main                    |
      | main    | git merge --squash feature           |
      |         | git commit -m "with "double quotes"" |
      |         | git push                             |
      |         | git push origin :feature             |
      |         | git branch -D feature                |
    And I am now on the "main" branch
    And the branches are now
      | REPOSITORY    | BRANCHES |
      | local, origin | main     |
    And my repo doesn't have any uncommitted files
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | with "double quotes" |
    And Git Town is now aware of no branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git push -u origin feature                    |
      |         | git revert {{ sha 'with "double quotes"' }}   |
      |         | git push                                      |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git checkout feature                          |
    And I am now on the "feature" branch
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE                       |
      | main    | local, origin | with "double quotes"          |
      |         |               | Revert "with "double quotes"" |
      | feature | local, origin | feature commit                |
    And the initial branch setup and hierarchy exists
