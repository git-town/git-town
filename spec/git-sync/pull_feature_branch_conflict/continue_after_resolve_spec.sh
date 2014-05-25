require_remote_main_branch

describe "git sync with conflicts after pulling the feature branch"

  context "when the user continues after resolving the conflicts"

    function before {
      create_feature_branch $feature_branch_name
      push_feature_branch
      add_local_commit $main_branch_name 'main_branch_update'
      add_remote_commit $feature_branch_name 'conflicting_remote_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_local_commit' 'conflicting_file' 'two'
      git sync
      git add conflicting_file
      git rebase --continue
      git sync --continue
    }

    it 'remains on the feature branch'
      expect_current_branch_is $feature_branch_name

    it 'keeps the local updates'
      expect_file_exists 'conflicting_file'
      expect_local_branch_has_commit $feature_branch_name 'conflicting_local_commit'

    it 'adds the resolved remote commit'
      expect_local_branch_has_commit $feature_branch_name 'conflicting_remote_commit'

    it 'continues with the sync, for example by adding updates from the main branch'
      expect_file_exists 'main_branch_update'

    it 'pushes the feature branch to the repo'
      expect_synchronized_branch $feature_branch_name

    it 'deletes the sync script'
      expect_file_does_not_exist "/tmp/git_sync_continue$temp_filename_suffix"

