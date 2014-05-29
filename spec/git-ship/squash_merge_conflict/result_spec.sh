require_remote_main_branch

describe "git-ship with conflicts after the squash merge of the feature branch into the main branch"

  context "situation after the conflict happens"

    function before {
      rm "/tmp/git_ship_abort$temp_filename_suffix"
      create_feature_branch $feature_branch_name
      add_local_commit $main_branch_name 'conflicting_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_commit' 'conflicting_file' 'two'
      git ship
    }

    it "ends up on the main branch"
      expect_current_branch_is $main_branch_name

    it "aborts in the middle of the merge"
      expect_conflict_for_file 'conflicting_file'

    it "creates an abort script file"
      expect_file_exists "/tmp/git_ship_abort$temp_filename_suffix"

    function after {
      git reset --hard HEAD
      git checkout $main_branch_name
      delete_feature_branch 'force'
    }

