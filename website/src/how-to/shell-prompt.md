# Display the currently pending Git Town command in your shell prompt

`git town status --pending` displays the name of the currently pending Git Town
command. This allows you to have your shell prompt display a reminder when you
need to run [git town continue](../commands/continue.md).

<img width="108" height="31" src="shell_prompt_example.gif">

### Bash

To add this status indicator to your Bash prompt, add the following to your
`.bashrc`:

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

For Zsh, add the following to your `~/.zshrc`:

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

For Fish shell, update your `~/.config/fish/config.fish` and override the
[`fish_prompt`](https://fishshell.com/docs/current/cmds/fish_prompt.html)
function:

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
