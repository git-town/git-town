# Tools that help create configuration files for testing.


config_path="$test_repo_path/.main_branch_name"


# Create a config file for the test directory
function create_test_config_file {
  echo "creating config file at $config_path with $main_branch_name"
  echo $main_branch_name > $config_path
}

