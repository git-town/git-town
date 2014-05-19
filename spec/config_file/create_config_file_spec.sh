describe "ensure_config_file_exists"

  context "when the config file exists"

    it "passes"
      ensure_config_file_exists


  context "when the config file does not exist"

    function before {
      # back up the config file
      mv $config_path "$config_path-old"

      # Run the SUT with mock user input
      echo "user_main_branch" > input
      ensure_config_file_exists < input
      rm input
    }

    function after {
      # Restore the config file
      mv "$config_path-old" $config_path
      main_branch_name=`cat $config_path`
    }

    it "creates a new config file with input queried from the user"
      expect_file_content $config_path "user_main_branch"

