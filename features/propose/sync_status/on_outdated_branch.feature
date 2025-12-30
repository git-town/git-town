@skipWindows
Feature: sync before proposing

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
    And the current branch is "child"
    And a proposal for this branch does not exist
    And tool "open" is installed
    When I run "git-town propose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                   |
      | child  | git fetch --prune --tags                                                  |
      |        | git checkout parent                                                       |
      | parent | git merge --no-edit --ff origin/parent                                    |
      |        | git push                                                                  |
      |        | git checkout child                                                        |
      | child  | git merge --no-edit --ff parent                                           |
      |        | git merge --no-edit --ff origin/child                                     |
      |        | git push                                                                  |
      |        | Finding proposal from child into parent ... ok                            |
      |        | open https://github.com/git-town/git-town/compare/parent...child?expand=1 |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE                                                  |
      | main   | local         | local main commit                                        |
      |        | origin        | origin main commit                                       |
      | parent | local, origin | local parent commit                                      |
      |        |               | origin parent commit                                     |
      |        |               | Merge remote-tracking branch 'origin/parent' into parent |
      | child  | local, origin | local child commit                                       |
      |        |               | Merge branch 'parent' into child                         |
      |        |               | origin child commit                                      |
      |        |               | Merge remote-tracking branch 'origin/child' into child   |
