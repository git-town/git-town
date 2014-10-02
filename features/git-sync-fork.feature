Feature: Git Sync-Fork

  Scenario: on the main branch with an upstream commit
    Given I am on the main branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | location | message         | file name     |
      | upstream | upstream commit | upstream_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`
    Then I am still on the "main" branch
    And I see the following commits
      | branch                  | location         | message         | files         |
      | main                    | local and remote | upstream commit | upstream_file |
      | remotes/upstream/main   | remote           | upstream commit | upstream_file |
    And now I have the following committed files
      | branch | files         |
      | main   | upstream_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"

  Scenario: on a feature branch with upstream commit in main branch
    Given I am on a feature branch
    And my repo has an upstream repo
    And the following commits exist in my repository
      | branch | location | message         | file name     |
      | main   | upstream | upstream commit | upstream_file |
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git sync-fork`
    Then I am still on the "feature" branch
    And I see the following commits
      | branch | location         | message         | files         |
      | main   | local and remote | upstream commit | upstream_file |
    And now I have the following committed files
      | branch | files         |
      | main   | upstream_file |
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"

  Scenario: without upstream configured and origin is not a GitHub url
    When I run `git sync-fork` while allowing errors
    Then I get the error "Please add a remote 'upstream'"

  @github_query
  Scenario: without upstream configured and origin is not found on GitHub
    Given my remote origin is "Originate/git-town-invalid" on GitHub
    When I run `git sync-fork` while allowing errors
    Then I get the error "Please add a remote 'upstream'"

  @github_query
  Scenario: without upstream configured and origin is not a fork on GitHub
    Given my remote origin is "Originate/git-town" on GitHub
    When I run `git sync-fork` while allowing errors
    Then I get the error "Please add a remote 'upstream'"

  @github_query
  Scenario: without upstream configured and origin is a fork on GitHub through HTTPS
    Given my remote origin is a "rails/rails" fork on GitHub through HTTPS
    When I run `git sync-fork --configure-only`
    Then my remote upstream is "rails/rails" on GitHub through HTTPS

  @github_query
  Scenario: without upstream configured and origin is a fork on GitHub through SSH
    Given my remote origin is a "rails/rails" fork on GitHub through SSH
    When I run `git sync-fork --configure-only`
    Then my remote upstream is "rails/rails" on GitHub through SSH
