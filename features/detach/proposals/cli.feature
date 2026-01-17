Feature: detach the current feature branch from a stack and update proposals

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | alpha  | feature | main   | local, origin |
      | beta   | feature | alpha  | local, origin |
      | gamma1 | feature | beta   | local, origin |
      | gamma2 | feature | beta   | local, origin |
      | delta  | feature | gamma2 | local, origin |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | TITLE           | BODY        | URL                      |
      |  1 | alpha         | main          | alpha proposal  | alpha body  | https://example.com/pr/1 |
      |  2 | beta          | alpha         | beta proposal   | beta body   | https://example.com/pr/2 |
      |  3 | gamma1        | beta          | gamma1 proposal | gamma1 body | https://example.com/pr/3 |
      |  4 | gamma2        | beta          | gamma2 proposal | gamma2 body | https://example.com/pr/4 |
      |  5 | delta         | gamma2        | delta proposal  | delta body  | https://example.com/pr/5 |
    And Git setting "git-town.proposals-show-lineage" is "cli"
    And the current branch is "beta"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                          |
      | beta   | git fetch --prune --tags                                         |
      |        | Finding proposal from beta into alpha ... #2 (beta proposal)     |
      |        | Finding proposal from gamma1 into beta ... #3 (gamma1 proposal)  |
      |        | Finding proposal from gamma2 into beta ... #4 (gamma2 proposal)  |
      |        | git checkout gamma1                                              |
      | gamma1 | git pull                                                         |
      |        | git -c rebase.updateRefs=false rebase --onto alpha beta          |
      |        | git push --force-with-lease                                      |
      |        | git checkout gamma2                                              |
      | gamma2 | git pull                                                         |
      |        | git -c rebase.updateRefs=false rebase --onto alpha beta          |
      |        | git push --force-with-lease                                      |
      |        | git checkout delta                                               |
      | delta  | git pull                                                         |
      |        | git -c rebase.updateRefs=false rebase --onto gamma2 beta         |
      |        | git push --force-with-lease                                      |
      |        | git checkout beta                                                |
      | beta   | git -c rebase.updateRefs=false rebase --onto main alpha          |
      |        | Updating target branch of proposal #2 to main ... ok             |
      |        | Updating target branch of proposal #3 to alpha ... ok            |
      |        | Updating target branch of proposal #4 to alpha ... ok            |
      |        | Finding all proposals for alpha ... main                         |
      |        | Finding proposal from gamma1 into alpha ... #3 (gamma1 proposal) |
      |        | Finding proposal from gamma2 into alpha ... #4 (gamma2 proposal) |
      |        | Finding proposal from delta into gamma2 ... #5 (delta proposal)  |
      |        | Update body for #1 ... ok                                        |
      |        | Finding all proposals for beta ... main                          |
      |        | Update body for #2 ... ok                                        |
      |        | Finding all proposals for delta ... gamma2                       |
      |        | Finding proposal from alpha into main ... #1 (alpha proposal)    |
      |        | Update body for #5 ... ok                                        |
      |        | Finding all proposals for gamma1 ... alpha                       |
      |        | Update body for #3 ... ok                                        |
      |        | Finding all proposals for gamma2 ... alpha                       |
      |        | Finding proposal from delta into gamma2 ... #5 (delta proposal)  |
      |        | Update body for #4 ... ok                                        |
    And this lineage exists now
      """
      main
        alpha
          gamma1
          gamma2
            delta
        beta
      """
    And the branches are now
      | REPOSITORY    | BRANCHES                                 |
      | local, origin | main, alpha, beta, delta, gamma1, gamma2 |
    And no uncommitted files exist now
    And the proposals are now
      """
      url: https://example.com/pr/1
      number: 1
      source: alpha
      target: main
      body:
        alpha body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1 :point_left:
            - https://example.com/pr/3
            - https://example.com/pr/4
              - https://example.com/pr/5
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/2
      number: 2
      source: beta
      target: main
      body:
        beta body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/2 :point_left:
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/3
      number: 3
      source: gamma1
      target: alpha
      body:
        gamma1 body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/3 :point_left:
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/4
      number: 4
      source: gamma2
      target: alpha
      body:
        gamma2 body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/4 :point_left:
              - https://example.com/pr/5
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/5
      number: 5
      source: delta
      target: gamma2
      body:
        delta body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/4
              - https://example.com/pr/5 :point_left:
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                               |
      |        | Updating target branch of proposal #2 to alpha ... ok |
      |        | Updating target branch of proposal #3 to beta ... ok  |
      |        | Updating target branch of proposal #4 to beta ... ok  |
    And the initial branches and lineage exist now
    And the initial commits exist now
    And the proposals are now
      """
      url: https://example.com/pr/1
      number: 1
      source: alpha
      target: main
      body:
        alpha body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1 :point_left:
            - https://example.com/pr/3
            - https://example.com/pr/4
              - https://example.com/pr/5
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/2
      number: 2
      source: beta
      target: alpha
      body:
        beta body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/2 :point_left:
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/3
      number: 3
      source: gamma1
      target: beta
      body:
        gamma1 body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/3 :point_left:
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/4
      number: 4
      source: gamma2
      target: beta
      body:
        gamma2 body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/4 :point_left:
              - https://example.com/pr/5
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      url: https://example.com/pr/5
      number: 5
      source: delta
      target: gamma2
      body:
        delta body
      
        <!-- branch-stack-start -->
      
        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/4
              - https://example.com/pr/5 :point_left:
      
        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>
      
        <!-- branch-stack-end -->
      
      """
