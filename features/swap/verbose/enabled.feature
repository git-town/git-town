Feature: swapping a feature branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT   | LOCATIONS     |
      | branch-1 | feature | main     | local, origin |
      | branch-2 | feature | branch-1 | local, origin |
      | branch-3 | feature | branch-2 | local, origin |
    And the commits
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | branch-1 | local, origin | commit 1    |
      | branch-2 | local, origin | commit 2    |
      | branch-3 | local, origin | commit 3    |
    And the current branch is "branch-2"
    When I run "git-town swap --verbose"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                               |
      |          | git version                                                                           |
      |          | git rev-parse --show-toplevel                                                         |
      |          | git config -lz --includes --global                                                    |
      |          | git config -lz --includes --local                                                     |
      |          | git status -z --ignore-submodules                                                     |
      |          | git rev-parse -q --verify MERGE_HEAD                                                  |
      |          | git rev-parse --absolute-git-dir                                                      |
      |          | git remote                                                                            |
      |          | git branch --show-current                                                             |
      | branch-2 | git fetch --prune --tags                                                              |
      | (none)   | git stash list                                                                        |
      |          | git -c core.abbrev=40 branch -vva --sort=refname                                      |
      |          | git remote get-url origin                                                             |
      |          | git rev-parse --verify --abbrev-ref @{-1}                                             |
      |          | git log --merges branch-1..branch-2                                                   |
      |          | git log --merges main..branch-1                                                       |
      | branch-2 | git -c rebase.updateRefs=false rebase --onto main branch-1                            |
      | (none)   | git rev-list --left-right branch-2...origin/branch-2                                  |
      | branch-2 | git push --force-with-lease --force-if-includes                                       |
      |          | git checkout branch-1                                                                 |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto branch-2 main                            |
      | (none)   | git rev-list --left-right branch-1...origin/branch-1                                  |
      | branch-1 | git push --force-with-lease --force-if-includes                                       |
      |          | git checkout branch-3                                                                 |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-1 {{ sha-before-run 'commit 2' }} |
      | (none)   | git rev-list --left-right branch-3...origin/branch-3                                  |
      | branch-3 | git push --force-with-lease --force-if-includes                                       |
      |          | git checkout branch-2                                                                 |
      | (none)   | git config git-town-branch.branch-2.parent main                                       |
      |          | git config git-town-branch.branch-1.parent branch-2                                   |
      |          | git config git-town-branch.branch-3.parent branch-1                                   |
      |          | git show-ref --verify --quiet refs/heads/branch-3                                     |
      |          | git -c core.abbrev=40 branch -vva --sort=refname                                      |
      |          | git config -lz --includes --global                                                    |
      |          | git config -lz --includes --local                                                     |
      |          | git stash list                                                                        |
    And Git Town prints:
      """
      Ran 36 shell commands.
      """

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                              |
      |          | git version                                          |
      |          | git rev-parse --show-toplevel                        |
      |          | git config -lz --includes --global                   |
      |          | git config -lz --includes --local                    |
      |          | git status -z --ignore-submodules                    |
      |          | git rev-parse -q --verify MERGE_HEAD                 |
      |          | git rev-parse --absolute-git-dir                     |
      |          | git stash list                                       |
      |          | git -c core.abbrev=40 branch -vva --sort=refname     |
      |          | git remote get-url origin                            |
      |          | git rev-parse --verify --abbrev-ref @{-1}            |
      |          | git remote get-url origin                            |
      | branch-2 | git checkout branch-1                                |
      | (none)   | git rev-parse HEAD                                   |
      | branch-1 | git reset --hard {{ sha 'commit 1' }}                |
      | (none)   | git rev-list --left-right branch-1...origin/branch-1 |
      | branch-1 | git push --force-with-lease --force-if-includes      |
      |          | git checkout branch-2                                |
      | (none)   | git rev-parse HEAD                                   |
      | branch-2 | git reset --hard {{ sha 'commit 2' }}                |
      | (none)   | git rev-list --left-right branch-2...origin/branch-2 |
      | branch-2 | git push --force-with-lease --force-if-includes      |
      |          | git checkout branch-3                                |
      | (none)   | git rev-parse HEAD                                   |
      | branch-3 | git reset --hard {{ sha 'commit 3' }}                |
      | (none)   | git rev-list --left-right branch-3...origin/branch-3 |
      | branch-3 | git push --force-with-lease --force-if-includes      |
      | (none)   | git show-ref --verify --quiet refs/heads/branch-2    |
      | branch-3 | git checkout branch-2                                |
      | (none)   | git config git-town-branch.branch-1.parent main      |
      |          | git config git-town-branch.branch-2.parent branch-1  |
      |          | git config git-town-branch.branch-3.parent branch-2  |
    And Git Town prints:
      """
      Ran 32 shell commands.
      """
