Feature: display all executed Git commands

  Background:
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git Town setting "sync-before-ship" is "true"

  Scenario: result
    When I run "git-town ship -m done --verbose"
    Then it runs the commands
      | BRANCH  | TYPE     | COMMAND                                           |
      |         | backend  | git version                                       |
      |         | backend  | git config -lz --global                           |
      |         | backend  | git config -lz --local                            |
      |         | backend  | git rev-parse --show-toplevel                     |
      |         | backend  | git status --long --ignore-submodules             |
      |         | backend  | git remote                                        |
      |         | backend  | git rev-parse --abbrev-ref HEAD                   |
      | feature | frontend | git fetch --prune --tags                          |
      |         | backend  | git stash list                                    |
      |         | backend  | git branch -vva --sort=refname                    |
      |         | backend  | git rev-parse --verify --abbrev-ref @{-1}         |
      |         | backend  | git remote get-url origin                         |
      | feature | frontend | git checkout main                                 |
      | main    | frontend | git rebase origin/main                            |
      |         | backend  | git rev-list --left-right main...origin/main      |
      | main    | frontend | git checkout feature                              |
      | feature | frontend | git merge --no-edit --ff origin/feature           |
      |         | frontend | git merge --no-edit --ff main                     |
      |         | backend  | git diff main..feature                            |
      | feature | frontend | git checkout main                                 |
      | main    | frontend | git merge --squash --ff feature                   |
      |         | backend  | git shortlog -s -n -e main..feature               |
      | main    | frontend | git commit -m done                                |
      |         | backend  | git rev-parse --short main                        |
      |         | backend  | git rev-list --left-right main...origin/main      |
      | main    | frontend | git push                                          |
      |         | frontend | git push origin :feature                          |
      |         | frontend | git branch -D feature                             |
      |         | backend  | git config --unset git-town-branch.feature.parent |
      |         | backend  | git show-ref --verify --quiet refs/heads/feature  |
      |         | backend  | git branch -vva --sort=refname                    |
      |         | backend  | git config -lz --global                           |
      |         | backend  | git config -lz --local                            |
      |         | backend  | git stash list                                    |
    And it prints:
      """
      Ran 34 shell commands.
      """
    And the current branch is now "main"

  Scenario: undo
    Given I ran "git-town ship -m done"
    When I run "git-town undo --verbose"
    Then it runs the commands
      | BRANCH | TYPE     | COMMAND                                        |
      |        | backend  | git version                                    |
      |        | backend  | git config -lz --global                        |
      |        | backend  | git config -lz --local                         |
      |        | backend  | git rev-parse --show-toplevel                  |
      |        | backend  | git status --long --ignore-submodules          |
      |        | backend  | git stash list                                 |
      |        | backend  | git branch -vva --sort=refname                 |
      |        | backend  | git rev-parse --verify --abbrev-ref @{-1}      |
      |        | backend  | git remote get-url origin                      |
      |        | backend  | git log --pretty=format:%h %s -10              |
      | main   | frontend | git revert {{ sha 'done' }}                    |
      |        | backend  | git rev-list --left-right main...origin/main   |
      | main   | frontend | git push                                       |
      |        | frontend | git branch feature {{ sha 'feature commit' }}  |
      |        | frontend | git push -u origin feature                     |
      |        | backend  | git show-ref --quiet refs/heads/feature        |
      | main   | frontend | git checkout feature                           |
      |        | backend  | git config git-town-branch.feature.parent main |
    And it prints:
      """
      Ran 18 shell commands.
      """
    And the current branch is now "feature"
