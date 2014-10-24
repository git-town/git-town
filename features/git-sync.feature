Feature: Git Sync

  Scenario: on the main branch
    Given I am on the main branch
    And the following commits exist in my repository
      | location | message       | file name   |
      | local    | local commit  | local_file  |
      | remote   | remote commit | remote_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "main" branch
    And all branches are now synchronized
    And I have the following commits
      | branch | location         | message       | files       |
      | main   | local and remote | local commit  | local_file  |
      | main   | local and remote | remote commit | remote_file |
    And now I have the following committed files
      | branch | files       |
      | main   | local_file  |
      | main   | remote_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: on a non feature branch
    Given non-feature branch configuration "qa, production"
    And I am on the "qa" branch
    And the following commits exist in my repository
      | branch | location         | message       | file name   |
      | qa     | local            | local commit  | local_file  |
      | qa     | remote           | remote commit | remote_file |
      | main   | local and remote | main commit   | main_file   |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "qa" branch
    And all branches are now synchronized
    And I have the following commits
      | branch | location         | message       | files       |
      | qa     | local and remote | local commit  | local_file  |
      | qa     | local and remote | remote commit | remote_file |
      | main   | local and remote | main commit   | main_file   |
    And now I have the following committed files
      | branch | files                   |
      | qa     | local_file, remote_file |
      | main   | main_file               |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: on a feature branch without a remote branch
    Given I am on a local feature branch
    And the following commits exist in my repository
      | branch  | location | message              | file name          |
      | main    | local    | local main commit    | local_main_file    |
      | main    | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "feature" branch
    And all branches are now synchronized
    And I have the following commits
      | branch  | location         | message                          | files               |
      | main    | local and remote | local main commit                | local_main_file     |
      | main    | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | Merge branch 'main' into feature |                     |
      | feature | local and remote | local main commit                | local_main_file     |
      | feature | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | local feature commit             | local_feature_file  |
    And now I have the following committed files
      | branch  | files                                                 |
      | main    | local_main_file, remote_main_file                     |
      | feature | local_feature_file, local_main_file, remote_main_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: on a feature branch with a remote branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message               | file name           |
      | main    | local    | local main commit     | local_main_file     |
      | main    | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      | feature | remote   | remote feature commit | remote_feature_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "feature" branch
    And I have the following commits
      | branch  | location         | message                          | files               |
      | main    | local and remote | local main commit                | local_main_file     |
      | main    | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | Merge branch 'main' into feature |                     |
      | feature | local and remote | local main commit                | local_main_file     |
      | feature | local and remote | remote main commit               | remote_main_file    |
      | feature | local and remote | local feature commit             | local_feature_file  |
      | feature | local and remote | remote feature commit            | remote_feature_file |
    And now I have the following committed files
      | branch  | files               |
      | main    | local_main_file, remote_main_file |
      | feature | local_feature_file, remote_feature_file, local_main_file, remote_main_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user aborts after a merge conflict when pulling the feature branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name          | file content               |
      | main    | local    | main branch update        | main_branch_update | main branch update         |
      | feature | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      | feature | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a merge in progress
    And there are abort and continue scripts for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location | message                   | files              |
      | main    | local    | main branch update        | main_branch_update |
      | feature | local    | local conflicting commit  | conflicting_file   |
      | feature | remote   | remote conflicting commit | conflicting_file   |
    And I still have the following committed files
      | branch  | files              | content                   |
      | feature | conflicting_file   | local conflicting content |
      | main    | main_branch_update | main branch update        |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user continues after resolving a merge conflict when pulling the feature branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name          | file content               |
      | feature | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      | feature | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
      | main    | local    | main branch update        | main_branch_update | main branch update         |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a merge in progress
    And there are abort and continue scripts for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I successfully finish the merge by resolving the merge conflict of file "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there are no abort and continue scripts for "git sync" anymore
    And now I have the following commits
      | branch  | location         | message                          | files              |
      | feature | local and remote | Merge branch 'main' into feature |                    |
      | feature | local and remote | remote conflicting commit        | conflicting_file   |
      | feature | local and remote | local conflicting commit         | conflicting_file   |
      | feature | local and remote | main branch update               | main_branch_update |
      | main    | local and remote | main branch update               | main_branch_update |
    And now I have the following committed files
      | branch  | files              | content            |
      | feature | conflicting_file   | resolved content   |
      | feature | main_branch_update | main branch update |
      | main    | main_branch_update | main branch update |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user aborts after a merge conflict when pulling the main branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch | location | message                   | file name          | file content               |
      | main   | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      | main   | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a rebase in progress
    And there are abort and continue scripts for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no rebase in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch | location | message                   | files              |
      | main   | local    | local conflicting commit  | conflicting_file   |
      | main   | remote   | remote conflicting commit | conflicting_file   |
    And I still have the following committed files
      | branch | files              | content                   |
      | main   | conflicting_file   | local conflicting content |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user aborts after a conflict when merging the main branch into the feature branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                   | file name        | file content    |
      | main    | local    | conflicting main commit   | conflicting_file | main content    |
      | feature | local    | conflicting local commit  | conflicting_file | feature content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a merge in progress
    And there are abort and continue scripts for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no merge in progress
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location | message                  | files            |
      | main    | local    | conflicting main commit  | conflicting_file |
      | feature | local    | conflicting local commit | conflicting_file |
    And I still have the following committed files
      | branch  | files            | content         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user continues after resolving the conflict when merging the main branch into the feature branch
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                    | file name        | file content    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a merge in progress
    And there are abort and continue scripts for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I successfully finish the merge by resolving the merge conflict of file "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there are no abort and continue scripts for "git sync" anymore
    And I still have the following commits
      | branch  | location         | message                          | files            |
      | main    | local            | conflicting main commit          | conflicting_file |
      | feature | local and remote | Merge branch 'main' into feature |                  |
      | feature | local and remote | conflicting main commit          | conflicting_file |
      | feature | local and remote | conflicting feature commit       | conflicting_file |
    And I still have the following committed files
      | branch  | files            | content          |
      | main    | conflicting_file | main content     |
      | feature | conflicting_file | resolved content |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  @finishes-with-non-empty-stash
  Scenario: user tries to continue without resolving an occurring merge conflict first
    Given I am on a feature branch
    And the following commits exist in my repository
      | branch  | location | message                    | file name        | file content    |
      | main    | local    | conflicting main commit    | conflicting_file | main content    |
      | feature | local    | conflicting feature commit | conflicting_file | feature content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a merge in progress
    When I run `git sync --continue` while allowing errors
    Then my repo still has a merge in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: two people collaborate on a feature branch
    Given I am on a feature branch
    And my coworker Charlie works on the same feature branch
    And the following commits exist in my repository
      | location  | message     | file name |
      | local     | my commit 1 | my_file_1 |
    And the following commits exist in Charlie's repository
      | location | message           | file name      |
      | local    | charlies commit 1 | charlie_file_1 |
    When I run `git sync`
    Then I see the following commits
      | branch  | location         | message     | files     |
      | feature | local and remote | my commit 1 | my_file_1 |
    And Charlie still sees the following commits
      | branch  | location | message           | files          |
      | feature | local    | charlies commit 1 | charlie_file_1 |
    When Charlie runs `git sync`
    Then now Charlie sees the following commits
      | branch  | location         | message           | files          |
      | feature | local and remote | charlies commit 1 | charlie_file_1 |
      | feature | local and remote | my commit 1       | my_file_1      |
    When I run `git sync`
    Then now I see the following commits
      | branch  | location         | message           | files          |
      | feature | local and remote | my commit 1       | my_file_1      |
      | feature | local and remote | charlies commit 1 | charlie_file_1 |


  Scenario: Unpushed tags on the main branch
    Given I am on the main branch
    And I add a local tag "v1.0"
    When I run `git sync`
    Then tag "v1.0" has been pushed to the remote


  Scenario: Unpushed tags on a non-feature branch
    Given non-feature branch configuration "qa, production"
    And I am on the "production" branch
    And I add a local tag "v1.0"
    When I run `git sync`
    Then tag "v1.0" has been pushed to the remote


  Scenario: On main branch with no pushable changes
    Given I am on the main branch
    And the following commits exist in my repository
      | location         | message               |
      | local and remote | already pushed commit |
    When I run `git sync`
    Then It doesn't run the command "git push"


  Scenario: On feature branch with no pushable changes
    Given I am on a feature branch
    And the following commits exist in my repository
      | location         | message               |
      | local and remote | already pushed commit |
    When I run `git sync`
    Then It doesn't run the command "git push"
