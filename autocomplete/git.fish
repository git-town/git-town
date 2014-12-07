# All Git Town commands
#
# TODO: create this variable programmatically, e.g.
# for command in (ls (brew --repository)/Library/Homebrew/cmd | sed -e "s/\.rb//g")
#   set commands $command $commands
# end
set commands extract hack kill pr prune-branches ship sync sync-fork town


# Indicates through its error code whether the command line to auto-complete
# already contains a Git Town command or not.
#
# - doesn't have command yet: exit code 0
# - has command already: exit code 1
function __fish_complete_git_town_no_command
  for cmd in (commandline -opc)
    if contains $cmd $commands
      return 1
    end
  end
  return 0
end


# define the Git Town commands
complete --command git --arguments 'extract' --description 'Copy selected commits from the current branch into their own branch' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'hack' --description 'Cut a new feature branch off the main branch' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'kill' --description 'Remove an obsolete feature branch' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'pr' --description 'Create a new pull request' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'prune-branches' --description 'Delete merged branches' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'ship' --description 'Deliver a completed feature branch' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'sync' --description 'Update the current branch with all relevant changes' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'sync-fork' --description 'Pull upstream updates into a forked repository' --condition '__fish_complete_git_town_no_command' --no-files
complete --command git --arguments 'town' --description 'Git Town management' --condition '__fish_complete_git_town_no_command' --no-files

# autocomplete Git branch names
complete --command git --arguments "(git branch | tr -d '* ')"

# command-line switches
complete --command git -l 'abort' --description 'Abort the current command'
complete --command git -l 'continue' --description 'Continue the current command'
complete --command git -l 'undo' --description 'Undo the current/last command'
