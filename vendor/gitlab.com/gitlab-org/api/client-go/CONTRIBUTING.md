# How to Contribute

We want to make contributing to this project as easy as possible.

## Reporting Issues

If you have an issue, please report it on the
[issue tracker](https://gitlab.com/gitlab-org/api/client-go/-/issues).

When you are up for writing a MR to solve the issue you encountered, it's not
needed to first open a separate issue. In that case only opening a MR with a
description of the issue you are trying to solve is just fine.

## Contributing Code

Merge requests are always welcome. When in doubt if your contribution fits within
the rest of the project, feel free to first open an issue to discuss your idea.

This is not needed when fixing a bug or adding an enhancement, as long as the
enhancement you are trying to add can be found in the public GitLab API docs as
this project only supports what is in the public API docs.

### Use community fork to contribute

To contribute to this project we recommend that you use the
[community fork](https://gitlab.com/gitlab-community/api/client-go).
Have a look at the
[community fork README](https://gitlab.com/gitlab-community#gitlab-community-forks)
to learn more about what it is and why you should prefer it over
creating your own fork to contribute.

## Coding style

We try to follow the Go best practices, where it makes sense, and use
[`gofumpt`](https://github.com/mvdan/gofumpt) to format code in this project.
As a general rule of thumb we prefer to keep line width for comments below 80
chars and for code (where possible and sensible) below 100 chars.

Before making a MR, please look at the rest this package and try to make sure
your contribution is consistent with the rest of the coding style.

New `struct` fields or methods should be placed (as much as possible) in the same
order as the ordering used in the public API docs. The idea is that this makes it
easier to find things.

### Setting up your local development environment to contribute

1. Install dependencies:
   ```sh
   make setup
   ```
1. Make your changes on your feature branch in the community fork or your personal fork
1. Run the reviewable command, which tests, lints and formats the code:
   ```sh
   make reviewable
   ```
1. Push your feature branch upstream
1. Open up your merge request
