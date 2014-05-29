require_remote_main_branch

describe "git-ship with conflicts after merging the feature branch into the main branch"

  context "when the user aborts git-ship"

    function before {
      rm "/tmp/git_ship_abort$temp_filename_suffix"
      create_feature_branch $feature_branch_name
      add_local_commit $main_branch_name 'conflicting_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_commit' 'conflicting_file' 'two'
      git ship
      git ship --abort
    }

    it "ends up on the feature branch"
      expect_current_branch_is $feature_branch_name

    it "aborts the merge"
      expected_no_merge_conflicts

    it "removes the abort script"
      expect_file_does_not_exist "/tmp/git_ship_abort$temp_filename_suffix"

    it "aborts the operation"
      expect_local_branch_exists $feature_branch_name


    function after {
      git rebase --abort
      git checkout $main_branch_name
      delete_feature_branch 'force'
    }

