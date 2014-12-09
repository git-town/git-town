#!/bin/bash -e
current_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

export program="$(echo "$0" | grep -o "[^/]*$")"

source "$current_dir/git_helpers/branch_helpers.sh"
source "$current_dir/git_helpers/checkout_helpers.sh"
source "$current_dir/git_helpers/cherry_pick_helpers.sh"
source "$current_dir/git_helpers/conflict_helpers.sh"
source "$current_dir/git_helpers/diff_helpers.sh"
source "$current_dir/git_helpers/feature_branch_helpers.sh"
source "$current_dir/git_helpers/merge_helpers.sh"
source "$current_dir/git_helpers/open_changes_helpers.sh"
source "$current_dir/git_helpers/rebase_helpers.sh"
source "$current_dir/git_helpers/remote_helpers.sh"
source "$current_dir/git_helpers/sha_helpers.sh"
source "$current_dir/git_helpers/tracking_branch_helpers.sh"

source "$current_dir/configuration.sh"
source "$current_dir/file_helpers.sh"
source "$current_dir/script_helpers.sh"
source "$current_dir/terminal_helpers.sh"
source "$current_dir/tool_helpers.sh"

export initial_branch_name=$(get_current_branch_name)
export initial_open_changes=$(has_open_changes)

if [ "$1" != "--abort" ] && [ "$1" != "--continue" ] && [ "$1" != "--undo" ]; then
  remove_scripts
fi
