# Collects error messages and information about which tests pass and which fail.

failures_log="$git_town_root_dir/errors"


function reset_failures_log {
  rm $failures_log 2> /dev/null
}
