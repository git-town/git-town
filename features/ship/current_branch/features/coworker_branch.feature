Feature: shipping a coworker's feature branch

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION      | MESSAGE         | AUTHOR                          |
      | feature | local, remote | coworker commit | coworker <coworker@example.com> |
    And I am on the "feature" branch

  Scenario: result (commit message via CLI)
    When I run "git-town ship -m 'feature done'"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                 |
      | feature | git fetch --prune --tags                                                |
      |         | git checkout main                                                       |
      | main    | git rebase origin/main                                                  |
      |         | git checkout feature                                                    |
      | feature | git merge --no-edit origin/feature                                      |
      |         | git merge --no-edit main                                                |
      |         | git checkout main                                                       |
      | main    | git merge --squash feature                                              |
      |         | git commit -m "feature done" --author "coworker <coworker@example.com>" |
      |         | git push                                                                |
      |         | git push origin :feature                                                |
      |         | git branch -D feature                                                   |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, remote | feature done | coworker <coworker@example.com> |
    And Git Town now has no branch hierarchy information

  Scenario: result (commit message via editor)
    When I run "git-town ship" and enter "feature done" for the commit message
    Then it runs the commands
      | BRANCH  | COMMAND                                               |
      | feature | git fetch --prune --tags                              |
      |         | git checkout main                                     |
      | main    | git rebase origin/main                                |
      |         | git checkout feature                                  |
      | feature | git merge --no-edit origin/feature                    |
      |         | git merge --no-edit main                              |
      |         | git checkout main                                     |
      | main    | git merge --squash feature                            |
      |         | git commit --author "coworker <coworker@example.com>" |
      |         | git push                                              |
      |         | git push origin :feature                              |
      |         | git branch -D feature                                 |
    And my repo now has the following commits
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, remote | feature done | coworker <coworker@example.com> |
    And Git Town now has no branch hierarchy information

  Scenario:  undo
    Given I ran "git-town ship -m 'feature done'"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                        |
      | main    | git branch feature {{ sha 'coworker commit' }} |
      |         | git push -u origin feature                     |
      |         | git revert {{ sha 'feature done' }}            |
      |         | git push                                       |
      |         | git checkout feature                           |
      | feature | git checkout main                              |
      | main    | git checkout feature                           |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, remote | feature done          |
      |         |               | Revert "feature done" |
      | feature | local, remote | coworker commit       |
    And Git Town now has the original branch hierarchy
