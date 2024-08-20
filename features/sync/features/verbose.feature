Feature: display all executed Git commands

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | local main commit     |
      |         | origin   | origin main commit    |
      | feature | local    | local feature commit  |
      |         | origin   | origin feature commit |

  Scenario: result
    When I run "git-town sync --verbose"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                                   |
      |         | backend  | git version                                               |
      |         | backend  | git rev-parse --show-toplevel                             |
      |         | backend  | git config -lz --includes --global                        |
      |         | backend  | git config -lz --includes --local                         |
      |         | backend  | git status --long --ignore-submodules                     |
      |         | backend  | git remote                                                |
      |         | backend  | git rev-parse --abbrev-ref HEAD                           |
      | feature | frontend | git fetch --prune --tags                                  |
      |         | backend  | git stash list                                            |
      |         | backend  | git branch -vva --sort=refname                            |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}                 |
      |         | backend  | git log main..feature --format=%h                         |
      |         | backend  | git log --format=%B -n 1 {{ sha 'local feature commit' }} |
      | feature | frontend | git checkout main                                         |
      | main    | frontend | git rebase origin/main                                    |
      |         | backend  | git rev-list --left-right main...origin/main              |
      | main    | frontend | git push                                                  |
      |         | frontend | git checkout feature                                      |
      | feature | frontend | git merge --no-edit --ff origin/feature                   |
      |         | frontend | git merge --no-edit --ff main                             |
      |         | backend  | git rev-list --left-right feature...origin/feature        |
      | feature | frontend | git push                                                  |
      |         | backend  | git show-ref --verify --quiet refs/heads/feature          |
      |         | backend  | git show-ref --verify --quiet refs/heads/main             |
      |         | backend  | git branch -vva --sort=refname                            |
      |         | backend  | git config -lz --includes --global                        |
      |         | backend  | git config -lz --includes --local                         |
      |         | backend  | git stash list                                            |
    And it prints:
      """
      Ran 28 shell commands.
      """
    And all branches are now synchronized
