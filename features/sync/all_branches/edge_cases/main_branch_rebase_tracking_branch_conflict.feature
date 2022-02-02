Feature: handling rebase conflicts between main branch and its tracking branch

  Background:
    Given my repo has a feature branch named "feature"
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE            | FILE NAME        | FILE CONTENT        |
      | main    | local    | main local commit  | conflicting_file | main local content  |
      |         | remote   | main remote commit | conflicting_file | main remote content |
      | feature | local    | feature commit     | feature_file     | feature content     |
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town sync --all"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git rebase origin/main   |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And my uncommitted file is stashed
    And my repo now has a rebase in progress

  Scenario: aborting
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH | COMMAND            |
      | main   | git rebase --abort |
      |        | git stash pop      |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo is left with my original commits

  Scenario: continuing without resolving the conflicts
    When I run "git-town continue"
    Then it runs no commands
    And it prints the error:
      """
      you must resolve the conflicts before continuing
      """
    And my uncommitted file is stashed
    And my repo still has a rebase in progress

  Scenario: continuing after resolving the conflicts
    Given I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git checkout main                  |
      | main    | git push --tags                    |
      |         | git stash pop                      |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
    And my repo now has the following commits
      | BRANCH  | LOCATION      | MESSAGE                          | FILE NAME        | FILE CONTENT        |
      | main    | local, remote | main remote commit               | conflicting_file | main remote content |
      |         |               | main local commit                | conflicting_file | resolved content    |
      | feature | local, remote | feature commit                   | feature_file     | feature content     |
      |         |               | main remote commit               | conflicting_file | main remote content |
      |         |               | main local commit                | conflicting_file | resolved content    |
      |         |               | Merge branch 'main' into feature |                  |                     |
    And my repo now has the following committed files
      | BRANCH  | NAME             | CONTENT          |
      | main    | conflicting_file | resolved content |
      | feature | conflicting_file | resolved content |
      |         | feature_file     | feature content  |

  Scenario: continuing after resolving the conflicts and continuing the rebase
    Given I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git push                           |
      |         | git checkout main                  |
      | main    | git push --tags                    |
      |         | git stash pop                      |
    And I am now on the "main" branch
    And my workspace has the uncommitted file again
