#!/usr/bin/env bash
# Note: "set -e" causes failures here

current_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export PROGRAM="$(echo "$0" | grep -o "[^/]*$")"
export GIT_COMMAND="${PROGRAM/-/ }"

bash "$current_dir/git_helpers/author_helpers.sh"
bash "$current_dir/git_helpers/branch_helpers.sh"
bash "$current_dir/git_helpers/branch_configuration_helpers.sh"
bash "$current_dir/git_helpers/checkout_helpers.sh"
bash "$current_dir/git_helpers/cherry_pick_helpers.sh"
bash "$current_dir/git_helpers/commits_helpers.sh"
bash "$current_dir/git_helpers/conflict_helpers.sh"
bash "$current_dir/git_helpers/fetch_helpers.sh"
bash "$current_dir/git_helpers/feature_branch_helpers.sh"
bash "$current_dir/git_helpers/merge_helpers.sh"
bash "$current_dir/git_helpers/open_changes_helpers.sh"
bash "$current_dir/git_helpers/push_helpers.sh"
bash "$current_dir/git_helpers/rebase_helpers.sh"
bash "$current_dir/git_helpers/remote_helpers.sh"
bash "$current_dir/git_helpers/sha_helpers.sh"
bash "$current_dir/git_helpers/shippable_changes_helpers.sh"
bash "$current_dir/git_helpers/tracking_branch_helpers.sh"

bash "$current_dir/browser_helpers.sh"
bash "$current_dir/configuration_helpers.sh"
bash "$current_dir/driver_helpers.sh"
bash "$current_dir/file_helpers.sh"
bash "$current_dir/folder_helpers.sh"
bash "$current_dir/script_helpers.sh"
bash "$current_dir/string_helpers.sh"
bash "$current_dir/sync_helpers.sh"
bash "$current_dir/terminal_helpers.sh"
bash "$current_dir/tool_helpers.sh"
bash "$current_dir/undo_helpers.sh"

bash "$current_dir/environment.sh" "$@"
bash "$current_dir/configuration.sh" "$@"

if [ "$(is_git_repository)" == true ]; then
  temp_filename_suffix="$(git_root | tr '/' '_')"
  export HAS_REMOTE=$(has_remote_url)
  export IN_SUB_FOLDER=$(is_in_git_sub_directory)
  export INITIAL_BRANCH_NAME=$(get_current_branch_name)
  export INITIAL_DIRECTORY=$(pwd)
  export INITIAL_OPEN_CHANGES=$(has_open_changes)
  export STEPS_FILE="/tmp/${PROGRAM}_${temp_filename_suffix}"
  export UNDO_STEPS_FILE="/tmp/${PROGRAM}_undo_${temp_filename_suffix}"
fi
