# Functions that run specs


# Runs all the specs
function run_all_specs {
  for spec_file in `find $git_town_root_dir/spec -name *_spec.sh`; do
    local spec_name=${spec_file/$git_town_root_dir\//}
    echo_header "RUNNING $spec_name"
    run_spec $spec_file
  done
}


# Runs the given spec.
function run_spec {
  source $1 2>> $failures_log
  reset_test_repo
}

