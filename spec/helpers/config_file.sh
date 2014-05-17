# Tools that help create configuration files for testing.


config_path="$test_repo_path/$config_filename"


# Create a config file for the test directory
function create_test_config_file {
  echo $main_branch_name > $config_path
}

