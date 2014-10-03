#!/bin/sh -e
current_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source $current_dir/terminal_helpers.sh
source $current_dir/configuration.sh
source $current_dir/file_helpers.sh
source $current_dir/git_helpers.sh
source $current_dir/github_helpers.sh

feature_branch_name=`get_current_branch_name`
current_branch_name=`get_current_branch_name`
determine_open_changes
