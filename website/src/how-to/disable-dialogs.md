# Disabling the dialogs

Git Town tries to be an ergonomic CLI application. Rather than just telling you
what is missing and how to enter it, it queries the needed information from you
right there when it needs it.

If you prefer a more conventional CLI application, you can set this environment
variable while running Git Town:

```bash
export TERM=dumb
```

This suppresses the dialogs.
