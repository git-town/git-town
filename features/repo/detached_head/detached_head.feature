Feature: open the repo in detached state

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS |
      | branch | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE  |
      | branch | local    | commit 1 |
      |        | local    | commit 2 |
    And the current branch is "branch"
    And tool "open" is installed
    And I ran "git checkout HEAD^"
    When I run "git-town repo"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH                     | TYPE     | COMMAND                                   |
      | {{ sha-short 'commit 1' }} | frontend | open https://github.com/git-town/git-town |
