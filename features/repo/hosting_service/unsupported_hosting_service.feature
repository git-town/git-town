Feature: unsupported forge type

  Scenario:
    Given a Git repo with origin
    When I run "git-town repo"
    Then Git Town prints the error:
      """
      unsupported forge type

      This command requires hosting on one of these services:
      * Bitbucket
      * Bitbucket Data Center
      * Forgejo
      * GitHub
      * GitLab
      * Gitea
      """
