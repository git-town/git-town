# Error handling

Sometimes Git Town commands encounter problems that require the human user to
make a decision. When this happens, the command stops and prints an error
message. When you have resolved the issue, you can either:

- run `git continue` to continue executing the interrupted command, starting
  with the operation that failed,
- run `git undo` to undo the Git Town command and go back to where you started.

You can also run `git undo` after a Git Town command finished to undo the
changes it made. Run `git town status` to see the status of the running Git Town
command and which Git Town commands you can run to continue or undo it.
