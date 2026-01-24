Feature: swap a feature branch and update proposals

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
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
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | TITLE             | BODY          | URL                      |
      | 1  | branch-1      | main          | branch-1 proposal | branch-1 body | https://example.com/pr/1 |
      | 2  | branch-2      | branch-1      | branch-2 proposal | branch-2 body | https://example.com/pr/2 |
      | 3  | branch-3      | branch-2      | branch-3 proposal | branch-3 body | https://example.com/pr/3 |
    And Git setting "git-town.proposals-show-lineage" is "cli"
    And the current branch is "branch-2"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                                            |
      | branch-2 | git fetch --prune --tags                                                           |
      |          | Finding proposal from branch-3 into branch-2 ... #3 (branch-3 proposal)            |
      |          | Finding proposal from branch-2 into branch-1 ... #2 (branch-2 proposal)            |
      |          | Finding proposal from branch-1 into main ... #1 (branch-1 proposal)                |
      |          | Updating target branch of proposal #2 to main ... ok                               |
      |          | git -c rebase.updateRefs=false rebase --onto main branch-1                         |
      |          | git push --force-with-lease --force-if-includes                                    |
      |          | Updating target branch of proposal #1 to branch-2 ... ok                           |
      |          | git checkout branch-1                                                              |
      | branch-1 | git -c rebase.updateRefs=false rebase --onto branch-2 main                         |
      |          | git push --force-with-lease --force-if-includes                                    |
      |          | Updating target branch of proposal #3 to branch-1 ... ok                           |
      |          | git checkout branch-3                                                              |
      | branch-3 | git -c rebase.updateRefs=false rebase --onto branch-1 {{ sha-initial 'commit 2' }} |
      |          | git push --force-with-lease --force-if-includes                                    |
      |          | git checkout branch-2                                                              |
      |          | Finding all proposals for branch-1 ... branch-2                                    |
      |          | Finding proposal from branch-2 into main ... #2 (branch-2 proposal)                |
      |          | Finding proposal from branch-1 into branch-2 ... #1 (branch-1 proposal)            |
      |          | Finding proposal from branch-3 into branch-1 ... #3 (branch-3 proposal)            |
      |          | Update body for #1 ... ok                                                          |
      |          | Finding all proposals for branch-2 ... main                                        |
      |          | Finding proposal from branch-1 into branch-2 ... #1 (branch-1 proposal)            |
      |          | Update body for #2 ... ok                                                          |
      |          | Finding all proposals for branch-3 ... branch-1                                    |
      |          | Finding proposal from branch-2 into main ... #2 (branch-2 proposal)                |
      |          | Update body for #3 ... ok                                                          |
    And this lineage exists now
      """
      main
        branch-2
          branch-1
            branch-3
      """
    And these commits exist now
      | BRANCH   | LOCATION      | MESSAGE     |
      | main     | local, origin | main commit |
      | branch-2 | local, origin | commit 2    |
      | branch-1 | local, origin | commit 1    |
      | branch-3 | local, origin | commit 3    |
    And the proposals are now
      """
      url: https://example.com/pr/1
      number: 1
      source: branch-1
      target: branch-2
      body:
        branch-1 body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/2
            - https://example.com/pr/1 :point_left:
              - https://example.com/pr/3

        <sup>[Stack](https://www.git-town.com/how-to/proposal-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/2
      number: 2
      source: branch-2
      target: main
      body:
        branch-2 body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/2 :point_left:
            - https://example.com/pr/1
              - https://example.com/pr/3

        <sup>[Stack](https://www.git-town.com/how-to/proposal-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/3
      number: 3
      source: branch-3
      target: branch-1
      body:
        branch-3 body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/2
            - https://example.com/pr/1
              - https://example.com/pr/3 :point_left:

        <sup>[Stack](https://www.git-town.com/how-to/proposal-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH   | COMMAND                                                  |
      | branch-2 | git checkout branch-1                                    |
      | branch-1 | git reset --hard {{ sha 'commit 1' }}                    |
      |          | git push --force-with-lease --force-if-includes          |
      |          | git checkout branch-2                                    |
      | branch-2 | git reset --hard {{ sha 'commit 2' }}                    |
      |          | git push --force-with-lease --force-if-includes          |
      |          | git checkout branch-3                                    |
      | branch-3 | git reset --hard {{ sha 'commit 3' }}                    |
      |          | git push --force-with-lease --force-if-includes          |
      |          | Updating target branch of proposal #2 to branch-1 ... ok |
      |          | Updating target branch of proposal #1 to main ... ok     |
      |          | Updating target branch of proposal #3 to branch-2 ... ok |
      |          | git checkout branch-2                                    |
    And the initial lineage exists now
    And the initial commits exist now
    And the proposals are now
      """
      url: https://example.com/pr/1
      number: 1
      source: branch-1
      target: main
      body:
        branch-1 body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/2
            - https://example.com/pr/1 :point_left:
              - https://example.com/pr/3

        <sup>[Stack](https://www.git-town.com/how-to/proposal-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/2
      number: 2
      source: branch-2
      target: branch-1
      body:
        branch-2 body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/2 :point_left:
            - https://example.com/pr/1
              - https://example.com/pr/3

        <sup>[Stack](https://www.git-town.com/how-to/proposal-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/3
      number: 3
      source: branch-3
      target: branch-2
      body:
        branch-3 body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/2
            - https://example.com/pr/1
              - https://example.com/pr/3 :point_left:

        <sup>[Stack](https://www.git-town.com/how-to/proposal-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->
      """
