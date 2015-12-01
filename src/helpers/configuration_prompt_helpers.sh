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

  if [ "$header_shown" = false ]; then
    echo_configuration_header
    header_shown=true
  fi

  local main_branch_input

  while [ -z "$main_branch_input" ]; do
    if [ "$(is_main_branch_configured)" = true ]; then
      echo -n "Please specify the main development branch by name or number (current value: ${MAIN_BRANCH_NAME}): "
    else
      echo -n "Please specify the main development branch by name or number (current value: None): "
    fi

    read user_input
    if [[ $user_input =~ $numerical_regex ]] ; then
      main_branch_input="$(get_numbered_branch "$user_input")"
      if [ -z "$main_branch_input" ]; then
        echo_error_header
        echo_error "invalid branch number"
      fi
    elif [ -z "$user_input" ]; then
      if [ "$(is_main_branch_configured)" = true ]; then
        main_branch_input=$MAIN_BRANCH_NAME
      else
        echo_error_header
        echo_error "no input received"
      fi
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


  local perennial_branches_input=''

  while true; do
    echo -n "Please specify a perennial branch by name or number. Leave it blank to finish (current value: ${PERENNIAL_BRANCH_NAMES}): "      

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
}
