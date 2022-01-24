Feature: git town-ship: offline mode

  When offline
  I want to be able to ship branches on my local machine
  So that I can keep working as much as possible despite having no internet connection.

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
    And I am on the "feature" branch
    When I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git branch -D feature              |
    And I am now on the "main" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION | MESSAGE        |
      | main    | local    | feature done   |
      | feature | remote   | feature commit |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                       |
      | main    | git branch feature {{ sha 'feature commit' }} |
      |         | git revert {{ sha 'feature done' }}           |
      |         | git checkout feature                          |
      | feature | git checkout main                             |
      | main    | git reset --hard {{ sha 'Initial commit' }}   |
      |         | git checkout feature                          |
    And I am now on the "feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, remote | feature commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |
