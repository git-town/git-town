Feature: git town-hack: starting a new feature from the main branch (without remote repo)

  To work with local
  I want to create new up-to-date feature branches and bring over my work to them.

  Background:
    Given my repo does not have a remote origin
    And the following commits exist in my repo
      | BRANCH | LOCATION | MESSAGE     |
      | main   | local    | main_commit |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git add -A                  |
      |             | git stash                   |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I am now on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION | MESSAGE     |
      | main        | local    | main_commit |
      | new-feature | local    | main_commit |

  Scenario: undo
    When I run "git town undo"
    Then it runs the commands
      | BRANCH      | COMMAND                   |
      | new-feature | git add -A                |
      |             | git stash                 |
      |             | git checkout main         |
      | main        | git branch -d new-feature |
      |             | git stash pop             |
    And I am now on the "main" branch
    And my repo now has the following commits
      | BRANCH      | LOCATION | MESSAGE     |
      | main        | local    | main_commit |
    And Git Town now has no branch hierarchy information
