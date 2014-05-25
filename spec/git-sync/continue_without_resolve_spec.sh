require_remote_main_branch

describe "git sync with conflicts after pulling the feature branch"

  context "when the user continues after resolving the conflicts"

    function before {
      create_feature_branch $feature_branch_name
      push_feature_branch
      add_remote_commit $feature_branch_name 'conflicting_remote_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_local_commit' 'conflicting_file' 'two'
      git sync
      git sync --continue
    }

    it 'does nothing'
      expect_rebase_in_progress


    function after {
      git rebase --abort
    }
