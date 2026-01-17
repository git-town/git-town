Feature: sync stacked changes and update proposals

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | parent | feature | main   | local, origin |
      | child  | feature | parent | local, origin |
      | other  | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION | MESSAGE              |
      | main   | local    | local main commit    |
      |        | origin   | origin main commit   |
      | parent | local    | local parent commit  |
      |        | origin   | origin parent commit |
      | child  | local    | local child commit   |
      |        | origin   | origin child commit  |
      | other  | local    | local other commit   |
      |        | origin   | origin other commit  |
    And the proposals
      | ID | SOURCE BRANCH | TARGET BRANCH | TITLE           | BODY        | URL                      |
      | 1  | parent        | main          | parent proposal | parent body | https://example.com/pr/1 |
      | 2  | child         | parent        | child proposal  | child body  | https://example.com/pr/2 |
      | 3  | other         | main          | other proposal  | other body  | https://example.com/pr/3 |
    And Git setting "git-town.proposals-show-lineage" is "cli"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    And the current branch is "child"
    When I run "git-town sync --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                |
      | child  | git fetch --prune --tags                                                               |
      |        | git checkout main                                                                      |
      | main   | git -c rebase.updateRefs=false rebase origin/main                                      |
      |        | git push                                                                               |
      |        | git checkout parent                                                                    |
      | parent | git push --force-with-lease --force-if-includes                                        |
      |        | git -c rebase.updateRefs=false rebase origin/parent                                    |
      |        | git -c rebase.updateRefs=false rebase --onto main {{ sha-initial 'initial commit' }}   |
      |        | git push --force-with-lease --force-if-includes                                        |
      |        | git checkout child                                                                     |
      | child  | git push --force-with-lease --force-if-includes                                        |
      |        | git -c rebase.updateRefs=false rebase origin/child                                     |
      |        | git -c rebase.updateRefs=false rebase --onto parent {{ sha-initial 'initial commit' }} |
      |        | git push --force-with-lease --force-if-includes                                        |
      |        | Finding all proposals for child ... parent                                             |
      |        | Finding proposal from parent into main ... #1 (parent proposal)                        |
      |        | Finding proposal from child into parent ... #2 (child proposal)                        |
      |        | Update body for #2 ... ok                                                              |
      |        | Finding all proposals for other ... main                                               |
      |        | Finding proposal from other into main ... #3 (other proposal)                          |
      |        | Update body for #3 ... ok                                                              |
      |        | Finding all proposals for parent ... main                                              |
      |        | Finding proposal from child into parent ... #2 (child proposal)                        |
      |        | Update body for #1 ... ok                                                              |
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | other  | local         | local other commit   |
      |        | origin        | origin other commit  |
      | parent | local, origin | origin parent commit |
      |        |               | local parent commit  |
      | child  | local, origin | origin child commit  |
      |        |               | local child commit   |
    And the proposals are now
      """
      url: https://example.com/pr/1
      number: 1
      source: parent
      target: main
      body:
        parent body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/1 :point_left:
            - https://example.com/pr/2

        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/2
      number: 2
      source: child
      target: parent
      body:
        child body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/2 :point_left:

        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/3
      number: 3
      source: other
      target: main
      body:
        other body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/3 :point_left:

        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                                                                      |
      | child  | git reset --hard {{ sha-initial 'local child commit' }}                                      |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin child commit' }}:child   |
      |        | git checkout parent                                                                          |
      | parent | git reset --hard {{ sha-initial 'local parent commit' }}                                     |
      |        | git push --force-with-lease origin {{ sha-in-origin-initial 'origin parent commit' }}:parent |
      |        | git checkout child                                                                           |
    And the initial branches and lineage exist now
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE              |
      | main   | local, origin | origin main commit   |
      |        |               | local main commit    |
      | other  | local         | local other commit   |
      |        | origin        | origin other commit  |
      | parent | local         | local parent commit  |
      |        | origin        | origin parent commit |
      | child  | local         | local child commit   |
      |        | origin        | origin child commit  |
    And the proposals are now
      """
      url: https://example.com/pr/1
      number: 1
      source: parent
      target: main
      body:
        parent body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/1 :point_left:
            - https://example.com/pr/2

        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/2
      number: 2
      source: child
      target: parent
      body:
        child body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/1
            - https://example.com/pr/2 :point_left:

        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->

      url: https://example.com/pr/3
      number: 3
      source: other
      target: main
      body:
        other body

        <!-- branch-stack-start -->

        -------------------------
        - main
          - https://example.com/pr/3 :point_left:

        <sup>[Stack](https://www.git-town.com/how-to/github-actions-breadcrumb.html) generated by [Git Town](https://github.com/git-town/git-town)</sup>

        <!-- branch-stack-end -->
      """
