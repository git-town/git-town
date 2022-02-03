Feature: append in offline mode

  Background:
    Given Git Town is in offline mode
    And my repo has a feature branch "existing-feature"
    And my repo contains the commits
      | BRANCH           | LOCATION      | MESSAGE                 |
      | existing-feature | local, remote | existing feature commit |
    And I am on the "existing-feature" branch

  Scenario: result
    When I run "git-town append new-feature"
    Then it runs the commands
      | BRANCH           | COMMAND                                     |
      | existing-feature | git checkout main                           |
      | main             | git rebase origin/main                      |
      |                  | git checkout existing-feature               |
      | existing-feature | git merge --no-edit origin/existing-feature |
      |                  | git merge --no-edit main                    |
      |                  | git branch new-feature existing-feature     |
      |                  | git checkout new-feature                    |
    And I am now on the "new-feature" branch
    And my repo now has the commits
      | BRANCH           | LOCATION      | MESSAGE                 |
      | existing-feature | local, remote | existing feature commit |
      | new-feature      | local         | existing feature commit |

  Scenario: undo
    Given I ran "git-town append new-feature"
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH           | COMMAND                       |
      | new-feature      | git checkout existing-feature |
      | existing-feature | git branch -D new-feature     |
      |                  | git checkout main             |
      | main             | git checkout existing-feature |
    And I am now on the "existing-feature" branch
    And my repo is left with my original commits
    And Git Town now has the original branch hierarchy
