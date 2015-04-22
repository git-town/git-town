Feature: Allow checking out previous git branch to work correctly after running a Git Town commmand

  As a developer running `git checkout -` after running a Git Town command
  I want to end up on the expected previous branch
  So that I can consistently and effectively use git's commands


  Scenario Outline: Running Git Town commands that don't publicly switch branches
    Given I have feature branches named "previous"
    And my repo has an upstream repo
    And I am on the "previous" branch
    And I run `<COMMAND>`
    When I checkout the previous git branch
    Then I end up on the "previous" branch

    Examples:
      | COMMAND       |
      | git hack new  |
      | git sync      |
      | git sync-fork |
