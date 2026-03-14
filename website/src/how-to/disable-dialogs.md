# Disabling the interactive dialogs

Git Town is designed to be an ergonomic interactive command-line application.
When it needs additional input, it asks for it directly instead of failing and
telling you to re-run with additional flags.

If you prefer a traditional, non-interactive CLI workflow, you can disable these
dialogs by setting the terminal type to `dumb` when running Git Town:

```bash
export TERM=dumb
```

When this environment variable is set, Git Town suppresses all interactive
dialogs and behaves like a conventional CLI tool.
