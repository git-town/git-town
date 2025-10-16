Feature: proposing uncommitted changes via a separate top-level branch,  provide commit message via CLI flag

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And the current branch is "existing"
    And tool "open" is installed
    And an uncommitted file "new_file" with content "new content"
    And I ran "git add new_file"
    When I run "git-town hack new --propose -m unrelated"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                        |
      | existing | git checkout -b new main                                                       |
      | new      | git commit -m unrelated                                                        |
      |          | git push -u origin new                                                         |
      |          | open https://github.com/git-town/git-town/compare/new?expand=1&title=unrelated |
      |          | git checkout existing                                                          |
    And this lineage exists now
      """
      main
        existing
        new
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | main     | origin        | main commit     |
      | existing | local         | existing commit |
      | new      | local, origin | unrelated       |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND              |
      | existing | git branch -D new    |
      |          | git push origin :new |
    And the initial branches and lineage exist now
    And the initial commits exist now
