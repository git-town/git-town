Feature: do not undo branches that were pulled in through "git fetch" while resolving conflicts

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE    | PARENT | LOCATIONS     |
      | feature    | feature | main   | local, origin |
      | coworker-1 | feature | main   | origin        |
    And the commits
      | BRANCH     | LOCATION | MESSAGE                   | FILE NAME         | FILE CONTENT   |
      | feature    | local    | conflicting local commit  | conflicting_file  | local content  |
      |            | origin   | conflicting origin commit | conflicting_file  | origin content |
      | coworker-1 | origin   | coworker-1 commit A       | coworker_1_file_a | content 1A     |
    And the current branch is "feature"
    And a coworker clones the repository
    And I run "git-town sync"
    And Git Town runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git merge --no-edit --ff origin/feature |
    And Git Town prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And a merge is now in progress
    And I resolve the conflict in "conflicting_file"
    And the coworker pushes these commits to the "coworker-1" branch
      | MESSAGE             | FILE NAME         | FILE CONTENT |
      | coworker-1 commit B | coworker_1_file_b | content 1B   |
    And the coworker pushes a new "coworker-2" branch with these commits
      | MESSAGE             | FILE NAME         | FILE CONTENT |
      | coworker-2 commit A | coworker_2_file_a | content 2A   |
    And I run "git fetch"
    When I run "git-town continue"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND              |
      | feature | git commit --no-edit |
      |         | git push             |
    And no merge is now in progress
    And all branches are now synchronized

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                                                    |
      | feature | git reset --hard {{ sha 'conflicting local commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'conflicting origin commit' }}:feature |
    And no merge is now in progress
    And the branches are now
      | REPOSITORY | BRANCHES                              |
      | local      | main, feature                         |
      | origin     | main, coworker-1, coworker-2, feature |
    And these commits exist now
      | BRANCH     | LOCATION         | MESSAGE                   | FILE NAME         | FILE CONTENT   |
      | feature    | local            | conflicting local commit  | conflicting_file  | local content  |
      |            | origin           | conflicting origin commit | conflicting_file  | origin content |
      | coworker-1 | coworker, origin | coworker-1 commit A       | coworker_1_file_a | content 1A     |
      |            |                  | coworker-1 commit B       | coworker_1_file_b | content 1B     |
      | coworker-2 | coworker, origin | coworker-2 commit A       | coworker_2_file_a | content 2A     |
