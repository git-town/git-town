# Display the currently pending Git Town command in your shell prompt

You can display a reminder for running `git town continue` to finish a pending
Git Town command in your shell prompt. Here is what this could look like:

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
