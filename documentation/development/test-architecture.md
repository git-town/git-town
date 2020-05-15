# Git Town Test Architecture

In the Go-based Cucumber setup, a [GitManager](../../test/git_manager.go)
instance manages the various Git repositories needed for tests. For each
scenario, GitManager provides a standardized
[GitEnvironment](../../test/git_environment.go) that contains:

- a "developer" GitRepository with the local workspace repo. This is where the
  feature specs execute in.
- an "origin" GitRepository that acts as the remote repo for the "developer"
  repo. This is where pushes from "developer" go to.
- the root directory of the GitEnvironment acts as the HOME directory. It
  contains the global Git configuration to use in this test.

Setting up a GitEnvironment is an expensive operation and has to be done for
every scenario. As a performance optimization, the GitManager creates a fully
set up "memoized" environment (including "main" branch and configuration) as a
template and copies it into the folders of the scenarios.

When running Go-based Cucumber concurrently, all threads share a global
GitManager instance.
