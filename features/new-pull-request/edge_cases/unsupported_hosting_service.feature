@this
Feature: unsupported hosting service

  Background:
    Given the current branch is a feature branch "feature"
    When I run "git-town new-pull-request"

  Scenario: result
    Then it prints the error:
      """
      unsupported hosting service

      This command requires hosting on one of these services:
      * Bitbucket
      * GitHub
      * GitLab
      * Gitea
      * Azure DevOps
      """
