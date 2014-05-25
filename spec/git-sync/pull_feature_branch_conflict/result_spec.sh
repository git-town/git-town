require_remote_main_branch

describe "git sync with conflicts after pulling the feature branch"

  context "right after the conflict happens"

    function before {
      create_feature_branch $feature_branch_name
      push_feature_branch
      add_remote_commit $feature_branch_name 'conflicting_remote_commit' 'conflicting_file' 'one'
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

