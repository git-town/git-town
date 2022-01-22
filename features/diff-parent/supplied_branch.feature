Feature: Viewing changes made on another branch

  To know whether my global branch setup is correct
  When working with nested feature branches
  I want to see the changes that a particular feature branch makes.

  Scenario: feature branch with known parent
    Given my repo has a feature branch named "feature-1"
    And my repo has a feature branch named "feature-2" as a child of "feature-1"
    When I run "git-town diff-parent feature-2"
    Then it runs the commands
      | BRANCH | COMMAND                       |
      | main   | git diff feature-1..feature-2 |

  @skipWindows
  Scenario: feature branch with unknown parent
    Given my repo has a feature branch named "feature" with no parent
    And I am on the "main" branch
    When I run "git-town diff-parent feature" and answer the prompts:
      | PROMPT                                        | ANSWER  |
      | Please specify the parent branch of 'feature' | [ENTER] |
    Then it runs the commands
      | BRANCH | COMMAND                |
      | main   | git diff main..feature |
    And Git Town is now aware of this branch hierarchy
      | BRANCH  | PARENT |
      | feature | main   |

  Scenario: main branch
    When I run "git-town diff-parent main"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """

  Scenario: perennial branch
    And my repo has the perennial branch "qa"
    When I run "git-town diff-parent qa"
    Then it runs no commands
    And it prints the error:
      """
      you can only diff-parent feature branches
      """

  Scenario: non-existing branch
    When I run "git-town diff-parent non-existing"
    Then it runs no commands
    And it prints the error:
      """
      there is no local branch named "non-existing"
      """
