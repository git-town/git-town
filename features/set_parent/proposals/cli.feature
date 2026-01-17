Feature: change the parent of a feature branch and update proposals

  Background:
    Given a Git repo with origin
    And the branches
      | NAME       | TYPE    | PARENT | LOCATIONS     |
      | old-parent | feature | main   | local, origin |
      | new-parent | feature | main   | local, origin |
    And the commits
      | BRANCH     | LOCATION | MESSAGE           |
      | old-parent | local    | old parent commit |
      | new-parent | local    | new parent commit |
    And the branches
      | NAME  | TYPE    | PARENT     | LOCATIONS     |
      | child | feature | old-parent | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE      |
      | child  | local    | child commit |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | TITLE          | BODY       | URL                      |
      |  1 | child         | old-parent    | child proposal | child body | https://example.com/pr/1 |
    And Git setting "git-town.proposals-show-lineage" is "cli"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "child"
    When I run "git-town set-parent new-parent"

  @this
  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                            |
      | child  | git pull                                                           |
      |        | git -c rebase.updateRefs=false rebase --onto new-parent old-parent |
      |        | git push --force-with-lease --force-if-includes                    |
    And this lineage exists now
      """
      main
        new-parent
          child
        old-parent
      """
    And the proposals are now
      """
      url: https://example.com/pr/1
      number: 1
      source: child
      target: old-parent
      body:
        child body
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                |
      | child  | git reset --hard {{ sha 'child commit' }}                              |
      |        | git push --force-with-lease origin {{ sha 'old parent commit' }}:child |
    And the initial branches and lineage exist now
    And the initial commits exist now
