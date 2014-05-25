require_remote_main_branch

describe "git sync with conflicts after pulling the feature branch"

  context "when the user aborts git-sync"

    function before {
      create_feature_branch $feature_branch_name
      git br -a
      push_feature_branch
      add_remote_commit $feature_branch_name 'conflicting_remote_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_local_commit' 'conflicting_file' 'two'
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

