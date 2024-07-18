Feature: already existing remote branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | old      | feature | main   | local, origin |
      | existing | feature | main   | origin        |
    And the current branch is "old"
    And an uncommitted file
    When I run "git-town prepend existing"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                             |
      | old      | git add -A                          |
      |          | git stash                           |
      |          | git checkout main                   |
      | main     | git rebase origin/main              |
      |          | git checkout old                    |
      | old      | git merge --no-edit --ff origin/old |
      |          | git merge --no-edit --ff main       |
      |          | git checkout -b existing main       |
      | existing | git stash pop                       |
    And the current branch is now "existing"
    And the uncommitted file still exists
    And the initial commits exist
    And this lineage exists now
      | BRANCH   | PARENT   |
      | existing | main     |
      | old      | existing |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH   | COMMAND                |
      | existing | git add -A             |
      |          | git stash              |
      |          | git checkout old       |
      | old      | git branch -D existing |
      |          | git stash pop          |
    And the current branch is now "old"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial lineage exists
