require_remote_main_branch

describe "git sync"

  context "merge conflict at update_feature_branch"

    function before {
      create_feature_branch $feature_branch_name
      add_local_commit $main_branch_name 'conflicting_local_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_local_commit' 'conflicting_file' 'two'
      git sync
    }

    it "aborts in the middle of the rebase"
      expect_rebase_in_progress

    it "creates an .git_sync_abort file"
      expect_file_exists '/tmp/git_sync_abort'


    function after {
      git rebase --abort
      git checkout $main_branch_name
      delete_feature_branch 'force'
    }


  context "is aborted by the user after a merge conflict in pull_feature_branch"

    function before {
      create_feature_branch $feature_branch_name
      add_local_commit $main_branch_name 'conflicting_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_commit' 'conflicting_file' 'two'
      echo "temp" > open_file
      git sync
      git sync --abort
    }

    it "ends up on the feature branch"
      expect_current_branch_is $feature_branch_name

    it "aborts the rebase"
      expect_no_rebase_in_progress

    it "removes the abort script"
      expect_file_does_not_exist '/tmp/git_sync_abort'

    it "pops the stash"
      expect_file_exists open_file


    function after {
      git rebase --abort
      git checkout $main_branch_name
      delete_feature_branch 'force'
      git reset HEAD open_file
    }

