Feature: Git Sync

  Scenario: on the main branch
    Given I am on the main branch
    And the following commits exist
      | location | message       | file name   |
      | local    | local commit  | local_file  |
      | remote   | remote commit | remote_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "main" branch
    And all branches are now synchronized
    And I have the following commits
      | branch | message       | files       |
      | main   | local commit  | local_file  |
      | main   | remote commit | remote_file |
    And now I have the following committed files
      | branch | name        |
      | main   | local_file  |
      | main   | remote_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: on a feature branch without a remote branch
    Given I am on a local feature branch
    And the following commits exist
      | branch  | location | message              | file name          |
      | main    | local    | local main commit    | local_main_file    |
      | main    | remote   | remote main commit   | remote_main_file   |
      | feature | local    | local feature commit | local_feature_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "feature" branch
    And all branches are now synchronized
    And now I have the following committed files
      | branch  | name               |
      | main    | local_main_file    |
      | main    | remote_main_file   |
      | feature | local_feature_file |
      | feature | local_main_file    |
      | feature | remote_main_file   |
    And I have the following commits
      | branch  | message               | files               |
      | main    | local main commit     | local_main_file     |
      | main    | remote main commit    | remote_main_file    |
      | feature | local main commit     | local_main_file     |
      | feature | remote main commit    | remote_main_file    |
      | feature | local feature commit  | local_feature_file  |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: on a feature branch with a remote branch
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message               | file name           |
      | main    | local    | local main commit     | local_main_file     |
      | main    | remote   | remote main commit    | remote_main_file    |
      | feature | local    | local feature commit  | local_feature_file  |
      | feature | remote   | remote feature commit | remote_feature_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync`
    Then I am still on the "feature" branch
    And now I have the following committed files
      | branch  | name                |
      | main    | local_main_file     |
      | main    | remote_main_file    |
      | feature | local_feature_file  |
      | feature | remote_feature_file |
      | feature | local_main_file     |
      | feature | remote_main_file    |
    And I have the following commits
      | branch  | message               | files               |
      | main    | local main commit     | local_main_file     |
      | main    | remote main commit    | remote_main_file    |
      | feature | local main commit     | local_main_file     |
      | feature | remote main commit    | remote_main_file    |
      | feature | local feature commit  | local_feature_file  |
      | feature | remote feature commit | remote_feature_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user aborts after a merge conflict when pulling the feature branch
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message                   | file name          | file content               |
      | main    | local    | main branch update        | main_branch_update | main branch update         |
      | feature | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      | feature | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a rebase in progress
    And there is an abort script for "git sync"
    And there is a continue script for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no rebase in progress
    And there is no abort script for "git sync" anymore
    And there is no continue script for "git sync" anymore
    And I still have the following commits
      | branch  | message                   | files              |
      | main    | main branch update        | main_branch_update |
      | feature | local conflicting commit  | conflicting_file   |
    And I still have the following committed files
      | branch  | name               | content            |
      | feature | conflicting_file   | local conflicting content   |
      | main    | main_branch_update | main branch update |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user continues after resolving a merge conflict when pulling the feature branch
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message                   | file name          | file content               |
      | feature | remote   | remote conflicting commit | conflicting_file   | remote conflicting content |
      | feature | local    | local conflicting commit  | conflicting_file   | local conflicting content  |
      | main    | local    | main branch update        | main_branch_update | main branch update         |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a rebase in progress
    And there is an abort script for "git sync"
    And there is a continue script for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I successfully finish the rebase by resolving the merge conflict of file "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there is no abort script for "git sync" anymore
    And there is no continue script for "git sync" anymore
    And now I have the following commits
      | branch  | message                   | files              |
      | feature | remote conflicting commit | conflicting_file   |
      | feature | local conflicting commit  | conflicting_file   |
      | feature | main branch update        | main_branch_update |
      | main    | main branch update        | main_branch_update |
    And now I have the following committed files
      | branch  | name               | content            |
      | feature | conflicting_file   | resolved content   |
      | feature | main_branch_update | main branch update |
      | main    | main_branch_update | main branch update |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user aborts after a merge conflict when rebasing the feature branch against the main branch
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message                   | file name        | file content    |
      | main    | local    | conflicting main commit   | conflicting_file | main content    |
      | feature | local    | conflicting local commit  | conflicting_file | feature content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a rebase in progress
    And there is an abort script for "git sync"
    And there is a continue script for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I run `git sync --abort`
    Then I am still on the "feature" branch
    And there is no rebase in progress
    And there is no abort script for "git sync" anymore
    And there is no continue script for "git sync" anymore
    And I still have the following commits
      | branch  | message                   | files            |
      | main    | conflicting main commit   | conflicting_file |
      | feature | conflicting local commit  | conflicting_file |
    And I still have the following committed files
      | branch  | name             | content         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | feature content |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user continues after resolving a merge conflict when rebasing the feature branch against the main branch
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message                     | file name        | file content    |
      | main    | local    | conflicting main commit     | conflicting_file | main content    |
      | feature | local    | conflicting feature commit  | conflicting_file | feature content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a rebase in progress
    And there is an abort script for "git sync"
    And there is a continue script for "git sync"
    And I don't have an uncommitted file with name: "uncommitted"
    When I successfully finish the rebase by resolving the merge conflict of file "conflicting_file"
    And I run `git sync --continue`
    Then I am still on the "feature" branch
    And there is no abort script for "git sync" anymore
    And there is no continue script for "git sync" anymore
    And I still have the following commits
      | branch  | message                     | files            |
      | main    | conflicting main commit     | conflicting_file |
      | feature | conflicting main commit     | conflicting_file |
      | feature | conflicting feature commit  | conflicting_file |
    And I still have the following committed files
      | branch  | name             | content         |
      | main    | conflicting_file | main content    |
      | feature | conflicting_file | resolved content |
    And I again have an uncommitted file with name: "uncommitted" and content: "stuff"


  Scenario: user tries to continue without resolving an occurring merge conflict first
    Given I am on a feature branch
    And the following commits exist
      | branch  | location | message                     | file name        | file content    |
      | main    | local    | conflicting main commit     | conflicting_file | main content    |
      | feature | local    | conflicting feature commit  | conflicting_file | feature content |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync` while allowing errors
    Then my repo has a rebase in progress
    When I run `git sync --continue` while allowing errors
    Then my repo still has a rebase in progress
    And I don't have an uncommitted file with name: "uncommitted"


  Scenario: on a special branch
    Given I have branches named qa, production
    And I set special branch names to "qa, production"
    And the following commits exist
      | branch  | location | message            | file name       |
      | main    | local    | local main commit  | local_main_file |
    When I checkout the "qa" branch
    And I run `git sync` while allowing errors
    Then I get the error "qa is a special branch. Please checkout the main branch or a feature branch to sync"
    And I am still on the "qa" branch
    When I checkout the "production" branch
    And I run `git sync` while allowing errors
    Then I get the error "production is a special branch. Please checkout the main branch or a feature branch to sync"
    And I am still on the "production" branch
