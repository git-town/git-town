Feature: unsupported hosting service

  Scenario:
    When I run "git-town repo"
    Then it prints the error:
      """
      unsupported hosting service

      This command requires hosting on one of these services:
      * Bitbucket
      * GitHub
      * GitLab
      * Gitea
      """
