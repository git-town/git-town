Feature: unsupported hosting platform

  Scenario:
    Given a Git repo clone
    When I run "git-town repo"
    Then it prints the error:
      """
      unsupported hosting platform

      This command requires hosting on one of these services:
      * Bitbucket
      * GitHub
      * GitLab
      * Gitea
      """
