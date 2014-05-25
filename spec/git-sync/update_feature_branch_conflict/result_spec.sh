require_remote_main_branch

describe "git sync with a merge conflict when updating the feature branch with the main branch"

  context "right after the conflict"

    function before {
      create_feature_branch $feature_branch_name
      add_local_commit $main_branch_name 'conflicting_local_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_local_commit' 'conflicting_file' 'two'
      git sync
    }

    it "aborts in the middle of the rebase"
      expect_rebase_in_progress

    it "creates an abort script file"
      expect_file_exists "/tmp/git_sync_abort$temp_filename_suffix"

    it "creates an continue script file"
      expect_file_exists "/tmp/git_sync_continue$temp_filename_suffix"


    function after {
      git rebase --abort
      git checkout $main_branch_name
      delete_feature_branch 'force'
    }

