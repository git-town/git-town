Feature: proposing uncommitted changes via a separate top-level branch, let Git ask for the commit message

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | existing | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION | MESSAGE         |
      | main     | origin   | main commit     |
      | existing | local    | existing commit |
    And the current branch is "existing"
    And the origin is "git@github.com:git-town/git-town.git"
    And tool "open" is installed
    And an uncommitted file with name "new_file" and content "new content"
    And I ran "git add new_file"
    When I run "git-town hack new --propose" and enter "unrelated idea" for the commit message

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                        |
      | existing | git checkout -b new main                                       |
      | new      | git commit                                                     |
      |          | git push -u origin new                                         |
      | (none)   | open https://github.com/git-town/git-town/compare/new?expand=1 |
      | new      | git checkout existing                                          |
    And the current branch is still "existing"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE         |
      | main     | origin        | main commit     |
      | existing | local         | existing commit |
      | new      | local, origin | unrelated idea  |
    And this lineage exists now
      | BRANCH   | PARENT |
      | existing | main   |
      | new      | main   |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND              |
      | existing | git branch -D new    |
      |          | git push origin :new |
    And the current branch is now "existing"
    And the initial commits exist now
    And the initial branches and lineage exist now
