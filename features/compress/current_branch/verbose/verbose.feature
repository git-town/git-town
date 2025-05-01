Feature: compress the commits on a feature branch verbosely

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
      |         |               | commit 3 | file_3    | content 3    |
    And the current branch is "feature"
    When I run "git-town compress --verbose" and enter "compressed commit" for the commit message

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                                         |
      |         | git version                                                     |
      |         | git rev-parse --show-toplevel                                   |
      |         | git config -lz --includes --global                              |
      |         | git config -lz --includes --local                               |
      |         | git rev-parse --verify --abbrev-ref @{-1}                       |
      |         | git status -z --ignore-submodules                               |
      |         | git rev-parse --verify -q MERGE_HEAD                            |
      |         | git rev-parse --absolute-git-dir                                |
      |         | git remote                                                      |
      |         | git branch --show-current                                       |
      | feature | git fetch --prune --tags                                        |
      | (none)  | git stash list                                                  |
      |         | git -c core.abbrev=40 branch -vva --sort=refname                |
      |         | git remote get-url origin                                       |
      |         | git cherry -v main feature                                      |
      |         | git show --no-patch --format=%B {{ sha-before-run 'commit 1' }} |
      | feature | git reset --soft main                                           |
      |         | git commit -m "commit 1"                                        |
      | (none)  | git rev-list --left-right feature...origin/feature              |
      | feature | git push --force-with-lease --force-if-includes                 |
      | (none)  | git -c core.abbrev=40 branch -vva --sort=refname                |
      |         | git config -lz --includes --global                              |
      |         | git config -lz --includes --local                               |
      |         | git stash list                                                  |
    And Git Town prints:
      """
      Ran 24 shell commands
      """
    And all branches are now synchronized
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE  |
      | feature | local, origin | commit 1 |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"

  Scenario: undo
    When I run "git-town undo --verbose"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                                            |
      |         | git version                                        |
      |         | git rev-parse --show-toplevel                      |
      |         | git config -lz --includes --global                 |
      |         | git config -lz --includes --local                  |
      |         | git status -z --ignore-submodules                  |
      |         | git rev-parse --verify -q MERGE_HEAD               |
      |         | git rev-parse --absolute-git-dir                   |
      |         | git stash list                                     |
      |         | git -c core.abbrev=40 branch -vva --sort=refname   |
      |         | git remote get-url origin                          |
      |         | git rev-parse --verify --abbrev-ref @{-1}          |
      |         | git remote get-url origin                          |
      |         | git rev-parse HEAD                                 |
      | feature | git reset --hard {{ sha 'commit 3' }}              |
      | (none)  | git rev-list --left-right feature...origin/feature |
      | feature | git push --force-with-lease --force-if-includes    |
    And the initial commits exist now
    And the initial branches and lineage exist now
