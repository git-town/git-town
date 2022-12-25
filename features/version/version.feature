Feature: show the version of the current Git Town installation

  @this
  Scenario: outside a Git repository
    Given I am outside a Git repo
    When I run "git-town version"
    Then it prints:
      """
      Git Town v.*-dev
      """
