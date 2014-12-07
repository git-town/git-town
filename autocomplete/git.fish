# define the Git Town commands
complete --command git --arguments 'extract' --description 'Copy selected commits from the current branch into their own branch' --no-files
complete --command git --arguments 'hack' --description 'Cut a new feature branch off the main branch' --no-files
complete --command git --arguments 'kill' --description 'Remove an obsolete feature branch' --no-files
complete --command git --arguments 'pr' --description 'Create a new pull request' --no-files
complete --command git --arguments 'prune-branches' --description 'Delete merged branches' --no-files
complete --command git --arguments 'ship' --description 'Deliver a completed feature branch' --no-files
complete --command git --arguments 'sync' --description 'Update the current branch with all relevant changes' --no-files
complete --command git --arguments 'sync-fork' --description 'Pull upstream updates into a forked repository' --no-files
complete --command git --arguments 'town' --description 'Git Town ' --no-files

# command-line switches
complete --command git -l 'abort' --description 'Abort the current command'
complete --command git -l 'continue' --description 'Continue the current command'
complete --command git -l 'undo' --description 'Undo the current/last command'
