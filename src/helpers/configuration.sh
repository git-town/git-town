#!/usr/bin/env bash

# Performs initial configuration before running any Git Town command,
# unless the `--bypass-automatic-configuration` option is passed (used by git-town)


# Migrate old configuration (Git Town v0.2.2 and lower)
if [[ -f ".main_branch_name" ]]; then
  store_configuration main-branch-name "$(cat .main_branch_name)"
  rm .main_branch_name
fi


export MAIN_BRANCH_NAME=$(get_configuration main-branch-name)
export NON_FEATURE_BRANCH_NAMES=$(get_configuration non-feature-branch-names)


# Bypass the configuration if requested by caller (e.g. git-town)
if [[ "$@" =~ --bypass-automatic-configuration ]]; then
  return 0
fi

# Show configuration prompt
if [[ "$(is_git_town_configured)" == false ]]; then
  echo "Git Town hasn't been configured for this repository."
  echo "Please run 'git town config --setup'."
  echo -n "Would you like to do that now?"

  if prompt_yn; then
    echo
    setup_configuration
  else
    exit_with_abort
  fi
fi
