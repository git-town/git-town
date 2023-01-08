Feature: automatic upgrade of new-branch-push-flag to push-new-branches

  Scenario: old flag is set locally
    Given setting "new-branch-push-flag" is "true"
    When I run "git-town config"
    Then it prints:
      """
      I found the deprecated local setting "git-town.new-branch-push-flag".
      I am upgrading this setting to the new format "git-town.push-new-branches".
      """
    And setting "push-new-branches" is now "true"
    And setting "new-branch-push-flag" no longer exists locally

  Scenario: old flag is set globally
    Given setting "new-branch-push-flag" is globally "true"
    When I run "git-town config"
    Then it prints:
      """
      I found the deprecated global setting "git-town.new-branch-push-flag".
      I am upgrading this setting to the new format "git-town.push-new-branches".
      """
    And setting "push-new-branches" is now "true"
    And setting "new-branch-push-flag" no longer exists globally
