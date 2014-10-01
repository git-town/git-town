#!/bin/sh -e

current_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

source $current_dir/terminal_helpers.sh
source $current_dir/configuration.sh
source $current_dir/file_helpers.sh
source $current_dir/git_helpers.sh
source $current_dir/github_helpers.sh
