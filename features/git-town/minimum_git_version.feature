Feature: Enforce minimum Git version

  To set up Git Town correctly
  When having an unsupported version of Git installed
  I want to see guidance telling me the minimal required Git version.

  Reasoning: Using `git remote get-url <name>` which was added in 2.7.0
  https://github.com/git/git/blob/1eb437020a2c098a7c12da4c05082fbea10d98c9/Documentation/RelNotes/2.7.0.txt


  Scenario: using an unsupported Git Version
    Given my computer has Git "2.6.2" installed
    When I run "git-town config"
    Then it prints the error:
      """
      Git Town requires Git 2.7.0 or higher
      """
