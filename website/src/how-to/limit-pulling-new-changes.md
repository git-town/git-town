# Controlling when to pull new changes

By default, Git Town keeps your branches in sync to ensure your Git workspace
stays up to date with the work of your team members. But there are times when
you might want to avoid pulling the latest changes—like when you're reorganizing
your stack or trying to dodge costly rebuilds triggered by upstream updates.

In those cases, use
[git town sync --detached](../commands/sync.md#-d--detached--no-detached). It
syncs your stack without touching the main branch.
