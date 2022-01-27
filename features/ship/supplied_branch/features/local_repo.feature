Feature: shipping the supplied feature branch without a remote origin

  Background:
    Given my repo has the feature branches "feature" and "other-feature"
    And my repo does not have a remote origin
    And the following commits exist in my repo
      | BRANCH  | LOCATION | MESSAGE        |
      | feature | local    | feature commit |
    And I am on the "other-feature" branch
    And my workspace has an uncommitted file with name "feature_file" and content "conflicting content"
    When I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH        | COMMAND                      |
      | other-feature | git add -A                   |
      |               | git stash                    |
      |               | git checkout feature         |
      | feature       | git merge --no-edit main     |
      |               | git checkout main            |
      | main          | git merge --squash feature   |
      |               | git commit -m "feature done" |
      |               | git branch -D feature        |
      |               | git checkout other-feature   |
      | other-feature | git stash pop                |
    And I am now on the "other-feature" branch
    And my workspace still contains my uncommitted file
    And the existing branches are
      | REPOSITORY | BRANCHES            |
      | local      | main, other-feature |
    And my repo now has the following commits
      | BRANCH | LOCATION | MESSAGE      |
      | main   | local    | feature done |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | other-feature | main   |

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH        | COMMAND                                       |
      | other-feature | git add -A                                    |
      |               | git stash                                     |
      |               | git checkout main                             |
      | main          | git branch feature {{ sha 'feature commit' }} |
      |               | git revert {{ sha 'feature done' }}           |
      |               | git checkout feature                          |
      | feature       | git checkout other-feature                    |
      | other-feature | git stash pop                                 |
    And I am now on the "other-feature" branch
    And my repo now has the following commits
      | BRANCH  | LOCATION | MESSAGE               |
      | main    | local    | feature done          |
      |         |          | Revert "feature done" |
      | feature | local    | feature commit        |
    And Git Town is now aware of this branch hierarchy
      | BRANCH        | PARENT |
      | feature       | main   |
      | other-feature | main   |
