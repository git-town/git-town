Feature: git town-hack: resolving conflicts between uncommitted changes and the main branch

  Background:
    Given my repo has a feature branch named "existing-feature"
    And the following commits exist in my repo
      | BRANCH | LOCATION      | MESSAGE            | FILE NAME        | FILE CONTENT |
      | main   | local, remote | conflicting commit | conflicting_file | main content |
    And I am on the "existing-feature" branch
    And my workspace has an uncommitted file with name: "conflicting_file" and content: "conflicting content"
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
      The stash entry is kept in case you need it again.
      """
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And the file "conflicting_file" contains unresolved conflicts


  Scenario: aborting
    When I run "git-town abort"
  # TODO: make this work
  # Then it runs the commands
  # | BRANCH      | COMMAND           |
  # | new-feature | git checkout main |
  # And I end up on the "existing-feature" branch
  # And my workspace has the uncommitted file again
  # And my repo is left with my original commits


  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it prints the error:
      """
      you must resolve the conflicts before continuing
      """


  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    When I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH | COMMAND |
    And I end up on the "new-feature" branch
    And my workspace now contains the file "conflicting_file" with content "resolved content"
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE            | FILE NAME        |
      | main        | local, remote | conflicting commit | conflicting_file |
      | new-feature | local         | conflicting commit | conflicting_file |


  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    When I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH      | COMMAND                     |
      | main        | git push                    |
      |             | git branch new-feature main |
      |             | git checkout new-feature    |
      | new-feature | git stash pop               |
    And I end up on the "new-feature" branch
    And my workspace still contains my uncommitted file
    And my repo now has the following commits
      | BRANCH      | LOCATION      | MESSAGE                   | FILE NAME        |
      | main        | local, remote | conflicting remote commit | conflicting_file |
      |             |               | conflicting local commit  | conflicting_file |
      | new-feature | local         | conflicting remote commit | conflicting_file |
      |             |               | conflicting local commit  | conflicting_file |
