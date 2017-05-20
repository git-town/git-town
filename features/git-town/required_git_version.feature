Feature: git town: show an error message when minimum Git version is not satisfied

  As a user with an unsupported version of Git installed
  I want to be given a useful error messge
  So that I know my next step if I want to use Git Town


  Scenario: using an unsupported Git Version
    Given I have Git "2.5.6" installed
    When I run `git-town config`
    Then I get the error "Git Town requires Git 2.6.0 or higher"
