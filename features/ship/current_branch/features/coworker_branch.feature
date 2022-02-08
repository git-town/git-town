Feature: ship a coworker's feature branch

  Background:
    Given my repo has a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE         | AUTHOR                          |
      | feature | local, origin | coworker commit | coworker <coworker@example.com> |
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
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, origin | feature done | coworker <coworker@example.com> |
    And Git Town is now aware of no branch hierarchy

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
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE      | AUTHOR                          |
      | main   | local, origin | feature done | coworker <coworker@example.com> |
    And Git Town is now aware of no branch hierarchy

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
    And now these commits exist
      | BRANCH  | LOCATION      | MESSAGE               |
      | main    | local, origin | feature done          |
      |         |               | Revert "feature done" |
      | feature | local, origin | coworker commit       |
    And Git Town is now aware of the initial branch hierarchy
