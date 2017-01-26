#!/usr/bin/env bash


# Prints the header for the prompt when asking for configuration
function echo_configuration_header {
  echo "Git Town needs to be configured"
  echo
  echo_numbered_branches_alpha_order
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
  local main_branch_prompt='Please specify the main development branch by name or number'

  while [ -z "$main_branch_input" ]; do
    if [ "$(is_main_branch_configured)" = true ]; then
      echo -n "$main_branch_prompt (current value: $(echo_inline_cyan_bold "$MAIN_BRANCH_NAME")): "
    else
      echo -n "$main_branch_prompt (current value: None): "
    fi

    read user_input
    if [[ $user_input =~ $numerical_regex ]] ; then
      main_branch_input="$(get_numbered_branch_alpha_order "$user_input")"
      if [ -z "$main_branch_input" ]; then
        echo_error_header
        echo_error "Invalid branch number"
        echo
      fi
    elif [ -z "$user_input" ]; then
      if [ "$(is_main_branch_configured)" = true ]; then
        main_branch_input=$MAIN_BRANCH_NAME
      else
        echo_error_header
        echo_error "A main development branch is required to enable the features provided by Git Town"
        echo
      fi
    else
      if [ "$(has_branch "$user_input")" == true ]; then
        main_branch_input=$user_input
      else
        echo_error_header
        echo_error "Branch '$user_input' doesn't exist"
        echo
      fi
    fi
  done

  store_configuration main-branch-name "$main_branch_input"


  local perennial_branches_input=''
  local perennial_branches_prompt='Please specify a perennial branch by name or number. Leave it blank to finish'

  while true; do
    if [ "$(are_perennial_branches_configured)" = true ]; then
      echo -n "$perennial_branches_prompt (current value(s): $(echo_inline_cyan_bold "$PERENNIAL_BRANCH_NAMES")): "
    else
      echo -n "$perennial_branches_prompt (current value(s): None): "
    fi

    read user_input
    local branch
    if [[ $user_input =~ $numerical_regex ]] ; then
      branch="$(get_numbered_branch_alpha_order "$user_input")"
      if [ -z "$branch" ]; then
        echo_error_header
        echo_error "Invalid branch number"
        echo
      fi
    elif [ -z "$user_input" ]; then
      break
    else
      if [ "$(has_branch "$user_input")" == true ]; then
        if [ "$user_input" == "$MAIN_BRANCH_NAME" ]; then
          echo_error_header
          echo_error "'$user_input' is already set as the main branch"
          echo
        else
          branch=$user_input
        fi
      else
        echo_error_header
        echo_error "Branch '$user_input' doesn't exist"
        echo
      fi
    fi

    if [ -n "$branch" ]; then
      perennial_branches_input="$(insert_string "$perennial_branches_input" ' ' "$branch")"
    fi
  done

  store_configuration perennial-branch-names "$perennial_branches_input"
}
