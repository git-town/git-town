#!/bin/bash

# Helper methods for dealing with a browser


# Opens a browser to the given URL, errors if no tool is available
function open_browser {
  local url=$1

  local tools=(open xdg-open)
  for tool in "${tools[@]}"; do
    if [ "$(has_tool "$tool")" = true ]; then
      eval "$tool $url"
      return
    fi
  done

  echo_error_header
  echo_error "Opening a browser requires 'open' on Mac or 'xdg-open' on Linux."
  echo_error "If you would like another command to be supported,"
  echo_error "please open an issue at https://github.com/Originate/git-town/issues"
  exit_with_error newline
}
