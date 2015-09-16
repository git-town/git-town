#!/usr/bin/env bash


# Prints the header for the prompt when asking for configuration
function echo_configuration_header {
  echo "Git Town needs to be configured"
  echo
  echo_numbered_branches
  echo
}


# Makes sure Git Town is configured
function ensure_knows_configuration {
  local header_shown=false
  local numerical_regex='^[0-9]+$'
  local user_input

  if [ "$(is_main_branch_configured)" = false ]; then
    if [ "$header_shown" = false ]; then
      echo_configuration_header
      header_shown=true
    fi

    local main_branch_input

    while [ -z "$main_branch_input" ]; do
      echo -n "Please specify the main development branch by name or number: "
      read user_input
      if [[ $user_input =~ $numerical_regex ]] ; then
        main_branch_input="$(get_numbered_branch "$user_input")"
        if [ -z "$main_branch_input" ]; then
          echo_error_header
          echo_error "invalid branch number"
        fi
      elif [ -z "$user_input" ]; then
        echo_error_header
        echo_error "no input received"
      else
        if [ "$(has_branch "$user_input")" == true ]; then
          main_branch_input=$user_input
        else
          echo_error_header
          echo_error "branch '$user_input' doesn't exist"
        fi
      fi
    done

    store_configuration main-branch-name "$main_branch_input"
  fi

  if [ "$(are_perennial_branches_configured)" = false ]; then
    if [ "$header_shown" = false ]; then
      echo_configuration_header
      header_shown=true
    fi

    local perennial_branches_input=''

    while true; do
      echo -n "Please specify a perennial branch by name or number (blank line to finish): "
      read user_input
      local branch
      if [[ $user_input =~ $numerical_regex ]] ; then
        branch="$(get_numbered_branch "$user_input")"
        if [ -z "$branch" ]; then
          echo_error_header
          echo_error "invalid branch number"
        fi
      elif [ -z "$user_input" ]; then
        break
      else
        if [ "$(has_branch "$user_input")" == true ]; then
          if [ "$user_input" == "$MAIN_BRANCH_NAME" ]; then
            echo_error_header
            echo_error "'$user_input' is already set as the main branch"
          else
            branch=$user_input
          fi
        else
          echo_error_header
          echo_error "branch '$user_input' doesn't exist"
        fi
      fi

      if [ -n "$branch" ]; then
        perennial_branches_input="$(insert_string "$perennial_branches_input" ' ' "$branch")"
      fi
    done

    store_configuration perennial-branch-names "$perennial_branches_input"
  fi
}
