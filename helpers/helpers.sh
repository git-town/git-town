#!/bin/bash -e
current_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source "$current_dir/configuration.sh"
source "$current_dir/file_helpers.sh"
source "$current_dir/git_helpers.sh"
source "$current_dir/script_helpers.sh"
source "$current_dir/terminal_helpers.sh"

export initial_branch_name=$(get_current_branch_name)
export initial_open_changes=$(has_open_changes)
