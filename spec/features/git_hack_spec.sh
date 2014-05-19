require_remote_main_branch

describe "git hack"

  context "on the main branch"

    function before {
      git hack feature1
    }

    it "creates a new local branch with the given name"
      assert_local_branch_exists "feature1"

    it "checks out the new feature branch"
      assert_current_branch_is "feature1"

    it "does not push the new feature branch to the repo"
      assert_no_remote_branch_exists "feature1"


  context 'with uncommited local changes'

    function before {
      echo "content" > file1
    }

    it 'preserves the uncommitted changes'
      assert_uncommitted_changes "file1"

    function after {
      rm file1
    }


  context "on an outdated feature branch"

    function before {
      git checkout -b existing_feature
      add_local_commit $main_branch_name 'new_commit_in_main'
      checkout_branch 'existing_feature'
      git hack new_feature
    }

    it "cuts the new branch off the main branch"
      assert_branch_has_commit "new_commit_in_main"


  context "with updates from other developers on the main branch"

    function before {
      git checkout -b existing_feature
      add_remote_commit $main_branch_name 'new_remote_commit'
      checkout_branch 'existing_feature'
      git hack new_feature
    }

    it "updates the main branch with the latest changes from remote"
      checkout_branch $main_branch_name
      assert_branch_has_commit "new_remote_commit"
