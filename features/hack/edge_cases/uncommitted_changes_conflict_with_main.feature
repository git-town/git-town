Feature: conflicts between uncommitted changes and the main branch

  Background:
    Given my repo has a feature branch "existing-feature"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, remote | conflicting commit | conflicting_file | main content |
    And I am on the "existing-feature" branch
    And my workspace has an uncommitted file with name "conflicting_file" and content "conflicting content"
    When I run "git-town hack new-feature"

  Scenario: result
    Then it runs the commands
      | BRANCH           | COMMAND                     |
      | existing-feature | git fetch --prune --tags    |
      |                  | git add -A                  |
      |                  | git stash                   |
      |                  | git checkout main           |
      | main             | git rebase origin/main      |
      |                  | git branch new-feature main |
      |                  | git checkout new-feature    |
      | new-feature      | git stash pop               |
    And it prints the error:
      """
      conflicts between your uncommmitted changes and the main branch
      """
    And the file "conflicting_file" contains unresolved conflicts

  Scenario: abort without resolving the conflicts fails due to unresolved merge conflicts
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH      | COMMAND           |
      | new-feature | git checkout main |
    And it prints the error:
      """
      cannot check out branch "main"
      """
    And I am still on the "new-feature" branch

  Scenario: abort after resolving the conflicts undoes the hack and indicates the merge conflict in the uncommitted file
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH      | COMMAND                       |
      | new-feature | git checkout main             |
      | main        | git branch -d new-feature     |
      |             | git checkout existing-feature |
    And it prints the error:
      """
      cannot check out branch "existing-feature"
      """
    And I am now on the "main" branch
    And my repo is left with my original commits
    And my workspace now contains the file "conflicting_file" with content "resolved content"

  Scenario: continue without resolving the conflicts
    When I run "git-town continue"
    Then it prints the error:
      """
      you must resolve the conflicts before continuing
      """

  Scenario: continue after resolving the conflicts
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs no commands
    And I am now on the "new-feature" branch
    And my repo now has the commits
      | BRANCH      | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main        | local, remote | conflicting commit | conflicting_file | main content |
      | new-feature | local         | conflicting commit | conflicting_file | main content |
    And my workspace now contains the file "conflicting_file" with content "resolved content"

  Scenario: undo after resolving the conflicts undoes the hack but cannot get back to the original branch due to merge conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH      | COMMAND                       |
      | new-feature | git checkout main             |
      | main        | git branch -d new-feature     |
      |             | git checkout existing-feature |
    And it prints the error:
      """
      cannot check out branch "existing-feature"
      """
    And I am now on the "main" branch
    And my repo is left with my original commits
    And my workspace now contains the file "conflicting_file" with content "resolved content"
