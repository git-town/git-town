Feature: handle intermittent "git fetch" while resolving conflicts

  Background: I fetch updates while resolving merge conflicts
    Given a Git repo clone
    And the branches
      | NAME       | TYPE    | PARENT | LOCATIONS     |
      | feature    | feature | main   | local, origin |
      | coworker-1 | feature | main   | origin        |
    And the current branch is "feature"
    And the commits
      | BRANCH     | LOCATION | MESSAGE                   | FILE NAME         | FILE CONTENT   |
      | feature    | local    | conflicting local commit  | conflicting_file  | local content  |
      |            | origin   | conflicting origin commit | conflicting_file  | origin content |
      | coworker-1 | origin   | coworker-1 commit A       | coworker_1_file_a | content 1A     |
    And a coworker clones the repository
    And I run "git-town sync"
    And it runs the commands
      | BRANCH  | COMMAND                                 |
      | feature | git fetch --prune --tags                |
      |         | git checkout main                       |
      | main    | git rebase origin/main                  |
      |         | git checkout feature                    |
      | feature | git merge --no-edit --ff origin/feature |
    And it prints the error:
      """
      CONFLICT (add/add): Merge conflict in conflicting_file
      """
    And the current branch is still "feature"
    And a merge is now in progress
    And I resolve the conflict in "conflicting_file"
    And the coworker pushes these commits to the "coworker-1" branch
      | MESSAGE             | FILE NAME         | FILE CONTENT |
      | coworker-1 commit B | coworker_1_file_b | content 1B   |
    And the coworker pushes a new "coworker-2" branch with these commits
      | MESSAGE             | FILE NAME         | FILE CONTENT |
      | coworker-2 commit A | coworker_2_file_1 | content 2A   |
    And I run "git fetch"
    When I run "git-town continue"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                       |
      | feature | git commit --no-edit          |
      |         | git merge --no-edit --ff main |
      |         | git push                      |
    And no merge is in progress
    And all branches are now synchronized

  @this
  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                                                                    |
      | feature | git reset --hard {{ sha 'conflicting local commit' }}                                      |
      |         | git push --force-with-lease origin {{ sha-in-origin 'conflicting origin commit' }}:feature |
      |         | git push origin :coworker-2                                                                |
      |         | git push --force-with-lease origin {{ sha-in-origin 'coworker-1 commit A' }}:coworker-1    |
    And no merge is in progress
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH     | LOCATION         | MESSAGE                   | FILE NAME         | FILE CONTENT   |
      | coworker-1 | coworker, origin | coworker-1 commit A       | coworker_1_file_a | content 1A     |
      |            | coworker         | coworker-1 commit B       | coworker_1_file_b | content 1B     |
      |            |                  | coworker-2 commit A       | coworker_2_file_1 | content 2A     |
      | feature    | local            | conflicting local commit  | conflicting_file  | local content  |
      |            | origin           | conflicting origin commit | conflicting_file  | origin content |
    And the initial branches and lineage exist
