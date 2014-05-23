require_remote_main_branch

describe "git hack"

  context "on the main branch"

    function before {
      git hack $feature_branch_name
    }

    it "creates a new local branch with the given name"
      expect_local_branch_exists $feature_branch_name

    it "checks out the new feature branch"
      expect_current_branch_is $feature_branch_name

    it "does not push the new feature branch to the repo"
      expect_no_remote_branch_exists $feature_branch_name


  context 'with uncommited local changes'

    function before {
      echo "content" > file1
    }

    it 'preserves the uncommitted changes'
      expect_uncommitted_changes "file1"

    function after {
      rm file1
    }


  context "on an outdated feature branch"

    function before {
      git checkout -b existing_feature
      add_local_commit $main_branch_name 'new_local_commit_in_main'
      add_remote_commit $main_branch_name 'new_remote_commit_in_main'
      checkout_branch 'existing_feature'
      git hack new_feature
    }

    it "checks out the new feature branch"
      expect_current_branch_is 'new_feature'

    it "cuts the new branch off the main branch"
      expect_local_branch_has_commit "new_feature" "new_local_commit_in_main"
      expect_local_branch_has_commit "new_feature" "new_remote_commit_in_main"

