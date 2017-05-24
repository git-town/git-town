Feature: git town: show an error message when minimum Git version is not satisfied

  As a user with an unsupported version of Git installed
  I want to be given a useful error messge
  So that I know my next step if I want to use Git Town

  Reasoning:
    - Using `git remote get-url <name>` which was added in 2.7.0
      https://github.com/git/git/blob/1eb437020a2c098a7c12da4c05082fbea10d98c9/Documentation/RelNotes/2.7.0.txt


  Scenario: using an unsupported Git Version
    Given I have Git "2.6.2" installed
    When I run `git-town config`
    Then I get the error "Git Town requires Git 2.7.0 or higher"
