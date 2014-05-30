require_remote_main_branch

describe "git extract"

  context "on the main branch"

    function before {
      git extract $feature_branch_name
    }


    it 'remains on the main branch'
      expect_current_branch_is $main_branch_name

    it 'does nothing'
      expect_local_branch_count 2

