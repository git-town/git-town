Feature: deleting a remote only branch

  Background:
    Given my origin has a feature branch "feature"
    And my repo contains the commits
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | remote   | feature commit |
    And I am on the "main" branch
    And I run "git-town sync"
    When I run "git-town kill feature"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
      |        | git push origin :feature |
    And the existing branches are
      | REPOSITORY    | BRANCHES |
      | local, remote | main     |
    And Git Town now has no branch hierarchy information

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                                                 |
      | main   | git push origin {{ sha-in-remote 'feature commit' }}:refs/heads/feature |
    And my repo now has the initial branches
    And Git Town now has no branch hierarchy information
