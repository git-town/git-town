Feature: handle conflicts between the main branch and its tracking branch

  Background:
    Given my repo has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT    |
      | main    | local    | conflicting local commit  | conflicting_file | local content   |
      |         | remote   | conflicting remote commit | conflicting_file | remote content  |
      | feature | local    | feature commit            | feature_file     | feature content |
    And I am on the "feature" branch
    When I run "git-town ship -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git checkout main        |
      | main    | git rebase origin/main   |
    And it prints the error:
      """
      To abort, run "git-town abort".
      To continue after having resolved conflicts, run "git-town continue".
      """
    And my repo now has a rebase in progress

  Scenario: abort
    When I run "git-town abort"
    Then it runs the commands
      | BRANCH | COMMAND              |
      | main   | git rebase --abort   |
      |        | git checkout feature |
    And I am still on the "feature" branch
    And there is no rebase in progress anymore
    And my repo is left with my original commits
    And Git Town still has the original branch hierarchy

  Scenario: resolve the conflict and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git-town continue" and close the editor
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git rebase --continue              |
      |         | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
    And I am now on the "main" branch
    And the existing branches are
      | REPOSITORY    | BRANCHES |
      | local, remote | main     |
    And my repo now has the commits
      | BRANCH | LOCATION      | MESSAGE                   |
      | main   | local, remote | conflicting remote commit |
      |        |               | conflicting local commit  |
      |        |               | feature done              |
    And Git Town now has no branch hierarchy information

  Scenario: resolve the conflict, finish the rebase, and continue
    When I resolve the conflict in "conflicting_file"
    And I run "git rebase --continue" and close the editor
    And I run "git-town continue"
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | main    | git push                           |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit -m "feature done"       |
      |         | git push                           |
      |         | git push origin :feature           |
      |         | git branch -D feature              |
    And I am now on the "main" branch
