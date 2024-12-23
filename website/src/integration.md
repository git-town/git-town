# Integration

This page describes how to integrate Git Town into other applications.

## Git Town GitHub Action

The
[Git Town GitHub Action](https://github.com/marketplace/actions/git-town-github-action)
appends a visual breadcrumb of the currently reviewed branch and its location in
the its branch stack.

![example branch stack created by the Git Town GitHub Action](https://raw.githubusercontent.com/git-town/action/main/docs/example-visualization.png)

## Shell autocompletion

To have your shell auto-complete Git Town commands, set up
[shell autocompletion](commands/completions.md)

## Shell prompt

You can display a reminder for running `git town continue` to finish a pending
Git Town command in your shell prompt. Here is how this could look like:

<img width="108" height="31" src="shell_prompt_example.gif">

### Bash

To add the example status indicator to your shell prompt in Bash, add this to
your `.bashrc` file:

```bash
function git_town_status {
    local pending_gittown_command=$(git town status --pending)
    if [ -n "$pending_gittown_command" ]; then
      echo -e " \033[30;43m $pending_gittown_command \033[0m "
    fi
}

PS1='$(git_town_status)> '
```

### Zsh

To add the example status indicator to your shell prompt in Zsh, add this to
your `~/.zshrc` file:

```zsh
git_town_status() {
  local git_status
  git_status=$(git town status --pending)
  if [[ -n "$git_status" ]]; then
    echo "%K{yellow}%F{black} $git_status %f%k "
  fi
}

setopt PROMPT_SUBST
PROMPT='$(git_town_status)> '
```

### Fish

To add the example status indicator to your shell prompt in Fish, edit your
`~/.config/fish/config.fish` file and overwrite the
[`fish_prompt` function](https://fishshell.com/docs/current/cmds/fish_prompt.html):

```zsh
function fish_prompt
  set -f pending_gittown_command (git-town status --pending)
  if [ -n "$pending_gittown_command" ]
    set -f yellow_pending_gittown_command (set_color -b yellow)(set_color black)(echo " $pending_gittown_command ")(set_color normal)' '
  else
    set -f yellow_pending_gittown_command ''
  end
  printf '%s> ' $yellow_pending_gittown_command
end
```

## [Lazygit](https://github.com/jesseduffield/lazygit)

Example lazygit configuration file to integrate Git Town:

```yml
customCommands:
  - key: 'Y'
    context: 'global'
    description: 'Git-Town sYnc'
    command: 'git-town sync --all'
    stream: true
    loadingText: 'Syncing'
  - key: 'U'
    context: 'global'
    description: 'Git-Town Undo (undo the last git-town command)'
    command: 'git-town undo'
    prompts:
    - type: 'confirm'
      title: 'Undo Last Command'
      body: 'Are you sure you want to Undo the last git-town command?'
    stream: true
    loadingText: 'Undoing Git-Town Command'
  - key: '!'
    context: 'global'
    description: 'Git-Town Repo (opens the repo link)'
    command: 'git-town repo'
    stream: true
    loadingText: 'Opening Repo Link'
  - key: 'a'
    context: 'localBranches'
    description: "Git-Town Append"
    prompts:
      - type: 'input'
        title: "Enter name of new child branch. Branches off of '{{.CheckedOutBranch.Name}}'"
        key: 'BranchName'
    command: 'git-town append {{.Form.BranchName}}'
    stream: true
    loadingText: 'Appending'
  - key: 'h'
    context: 'localBranches'
    description: 'Git-Town Hack (creates a new branch)'
    prompts:
      - type: 'input'
        title: "Enter name of new branch. Branches off of 'Main'"
        key: 'BranchName'
    command: 'git-town hack {{.Form.BranchName}}'
    stream: true
    loadingText: 'Hacking'
  - key: 'K'
    context: 'localBranches'
    description: 'Git-Town Delete (deletes the current feature branch and sYnc)'
    command: 'git-town delete'
    prompts:
    - type: 'confirm'
      title: 'Delete current feature branch'
      body: 'Are you sure you want to delete the current feature branch?'
    stream: true
    loadingText: 'Deleting Feature Branch'
  - key: 'p'
    context: 'localBranches'
    description: 'Git-Town Propose (creates a pull request)'
    command: 'git-town propose'
    stream: true
    loadingText: 'Creating pull request'
  - key: 'P'
    context: 'localBranches'
    description: "Git-Town Prepend (creates a branch between the curent branch and its parent)"
    prompts:
      - type: 'input'
        title: "Enter name of the for child branch between '{{.CheckedOutBranch.Name}}' and its parent"
        key: 'BranchName'
    command: 'git-town prepend {{.Form.BranchName}}'
    stream: true
    loadingText: 'Prepending'
  - key: 'S'
    context: 'localBranches'
    description: 'Git-Town Skip (skip branch with merge conflicts when syncing)'
    command: 'git-town skip'
    stream: true
    loadingText: 'Skiping'
  - key: 'G'
    context: 'files'
    description: 'Git-Town GO aka:continue (continue after resolving merge conflicts)'
    command: 'git-town continue'
    stream: true
    loadingText: 'Continuing'
```
