Feature: commit message can contain double-quotes

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And I am on the "feature" branch
    When I run "git-town ship -m 'message containing "double quotes"'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                            |
      | feature | git fetch --prune --tags                           |
      |         | git checkout main                                  |
      | main    | git rebase origin/main                             |
      |         | git checkout feature                               |
      | feature | git merge --no-edit origin/feature                 |
      |         | git merge --no-edit main                           |
      |         | git checkout main                                  |
      | main    | git merge --squash feature                         |
      |         | git commit -m "message containing "double quotes"" |
      |         | git push                                           |
      |         | git push origin :feature                           |
      |         | git branch -D feature                              |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES |
      | local, remote | main     |
    And my repo doesn't have any uncommitted files
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE                            |
      | main   | local, remote | message containing "double quotes" |
    And Git Town now has no branch hierarchy information

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                   |
      | main    | git branch feature {{ sha 'feature commit' }}             |
      |         | git push -u origin feature                                |
      |         | git revert {{ sha 'message containing "double quotes"' }} |
      |         | git push                                                  |
      |         | git checkout feature                                      |
      | feature | git checkout main                                         |
      | main    | git checkout feature                                      |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                                     |
      | main    | local, remote | message containing "double quotes"          |
      |         |               | Revert "message containing "double quotes"" |
      | feature | local, remote | feature commit                              |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
