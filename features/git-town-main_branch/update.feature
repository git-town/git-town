Feature: set the main branch configuration

  As a user or tool configuring Git Town
  I want an easy way to specifically set the main branch
  So that I can configure Git Town safely, and the tool does exactly what I want.


  Scenario: main branch not yet configured
    Given I don't have a main branch name configured
    When I run `git-town main-branch main`
    Then Git Town prints no output
    And Git Town's main branch is now configured as "main"


  Scenario: main branch is configured
    Given my repository has the branches "main-old" and "main-new"
    And Git Town's main branch is configured as "main-old"
    When I run `git-town main-branch main-new`
    Then Git Town prints no output
    And Git Town's main branch is now configured as "main-new"


  Scenario: invalid branch name
    When I run `git-town main-branch non-existing`
    Then Git Town prints the error "no branch named 'non-existing'"
