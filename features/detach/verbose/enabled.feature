Feature: detaching an omni-branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | branch-1 | feature | main   | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-1 | local, origin | commit 1a |
      | branch-1 | local, origin | commit 1b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-2 | feature | branch-1 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-2 | local, origin | commit 2a |
      | branch-2 | local, origin | commit 2b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-3 | local, origin | commit 3a |
      | branch-3 | local, origin | commit 3b |
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-4 | feature | branch-3 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-4 | local, origin | commit 4a |
      | branch-4 | local, origin | commit 4b |
    And the current branch is "branch-2"
    When I run "git-town detach --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                              |
      |          | git version                                          |
      |          | git rev-parse --show-toplevel                        |
      |          | git config -lz --includes --global                   |
      |          | git config -lz --includes --local                    |
      |          | git status --long --ignore-submodules                |
      |          | git remote                                           |
      |          | git branch --show-current                            |
      | branch-2 | git fetch --prune --tags                             |
      | (none)   | git stash list                                       |
      |          | git branch -vva --sort=refname                       |
      |          | git remote get-url origin                            |
      |          | git rev-parse --verify --abbrev-ref @{-1}            |
      |          | git log --merges branch-1..branch-2                  |
      | branch-2 | git rebase --onto main branch-1                      |
      | (none)   | git rev-list --left-right branch-2...origin/branch-2 |
      | branch-2 | git push --force-with-lease --force-if-includes      |
      |          | git checkout branch-3                                |
      | branch-3 | git pull                                             |
      |          | git rebase --onto branch-1 branch-2                  |
      |          | git push --force-with-lease                          |
      | (none)   | git rev-list --left-right branch-3...origin/branch-3 |
      | branch-3 | git checkout branch-4                                |
      | branch-4 | git pull                                             |
      |          | git rebase --onto branch-3 branch-2                  |
      |          | git push --force-with-lease                          |
      | (none)   | git rev-list --left-right branch-4...origin/branch-4 |
      | branch-4 | git checkout branch-2                                |
      | (none)   | git config git-town-branch.branch-2.parent main      |
      |          | git config git-town-branch.branch-3.parent branch-1  |
      |          | git show-ref --verify --quiet refs/heads/main        |
      |          | git checkout main                                    |
      |          | git checkout branch-2                                |
      |          | git branch -vva --sort=refname                       |
      |          | git config -lz --includes --global                   |
      |          | git config -lz --includes --local                    |
      |          | git stash list                                       |
    And Git Town prints:
      """
      Ran 36 shell commands.
      """
    And the current branch is still "branch-2"
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE   |
      | branch-1 | local, origin | commit 1a |
      |          |               | commit 1b |
      | branch-2 | local, origin | commit 2a |
      |          |               | commit 2b |
      | branch-3 | local, origin | commit 3a |
      |          |               | commit 3b |
      | branch-4 | local, origin | commit 4a |
      |          |               | commit 4b |
    And this lineage exists now
      | BRANCH   | PARENT   |
      | branch-1 | main     |
      | branch-2 | main     |
      | branch-3 | branch-1 |
      | branch-4 | branch-3 |

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                         |
      | branch-2 | git reset --hard {{ sha 'commit 2b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-3                           |
      | branch-3 | git reset --hard {{ sha 'commit 3b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-4                           |
      | branch-4 | git reset --hard {{ sha 'commit 4b' }}          |
      |          | git push --force-with-lease --force-if-includes |
      |          | git checkout branch-2                           |
    And the current branch is still "branch-2"
    And the initial commits exist now
    And the initial lineage exists now
