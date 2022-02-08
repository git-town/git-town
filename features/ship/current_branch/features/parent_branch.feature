Feature: ship a parent branch

  Background:
    Given my repo has a feature branch "parent"
    And my repo has a feature branch "child" as a child of "parent"
    And my repo contains the commits
      | BRANCH | LOCATION      | MESSAGE       |
      | parent | local, origin | parent commit |
      | child  | local, origin | child commit  |
    And I am on the "parent" branch
    When I run "git-town ship -m 'parent done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                           |
      | parent | git fetch --prune --tags          |
      |        | git checkout main                 |
      | main   | git rebase origin/main            |
      |        | git checkout parent               |
      | parent | git merge --no-edit origin/parent |
      |        | git merge --no-edit main          |
      |        | git checkout main                 |
      | main   | git merge --squash parent         |
      |        | git commit -m "parent done"       |
      |        | git push                          |
      |        | git branch -D parent              |
    And I am now on the "main" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE       |
      | main   | local, origin | parent done   |
      | child  | local, origin | child commit  |
      | parent | origin        | parent commit |
    And Git Town is now aware of this branch hierarchy
      | BRANCH | PARENT |
      | child  | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                     |
      | main   | git branch parent {{ sha 'parent commit' }} |
      |        | git revert {{ sha 'parent done' }}          |
      |        | git push                                    |
      |        | git checkout parent                         |
      | parent | git checkout main                           |
      | main   | git checkout parent                         |
    And I am now on the "parent" branch
    And now these commits exist
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | parent done          |
      |        |               | Revert "parent done" |
      | child  | local, origin | child commit         |
      | parent | local, origin | parent commit        |
    And my repo now has its initial branches and branch hierarchy
