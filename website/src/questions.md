# Q & A

### Does this force me into any conventions for my branches or commits?

No. Git Town doesn’t require or enforce conventions for naming or set up of
branches and commits. It works with a wide variety of Git branching models and
workflows.

### Which Git branching models are supported by Git Town?

Git Town is so generic that it supports the most widely used branching models
including
[GitHub Flow](https://docs.github.com/en/get-started/quickstart/github-flow),
[Git Flow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow),
[GitLab Flow](https://docs.gitlab.com/ee/topics/gitlab_flow.html),
[trunk-based development](https://trunkbaseddevelopment.com) and even committing
straight into the main branch!

### How is this different from the [git-flow](https://github.com/nvie/gitflow) tool?

git-flow is a Git extension that provides specific and opinionated support for
the powerful Git branching model with the same name. It doesn’t care too much
about how you keep your work in sync with the rest of the team. Git Town doesn’t
care which branching model you use. It focusses on keeping your team
synchronized and your code repository clean. It is possible to use the two tools
together.

### Is Git Town compatible with my other Git tools?

Yes, we try to be good citizens in the Git ecosystem. If you run into any issues
with your setup, please let us know!

### Does my whole team have to use Git Town?

No. But please make sure that all feature branches get squash-merged, for
example by running `git merge --squash` or enabling them in your
[GitHub settings](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/configuring-commit-squashing-for-pull-requests).
If you don't know what squash-merges are, you probably want to learn about them
and use them.
