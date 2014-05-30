require_remote_main_branch

describe "git-extract with conflicts while cherry-picking"

  context "situation after the conflict happens"

    function before {
      new_branch_name="extracted_branch_$$_$current_time"
      create_feature_branch $feature_branch_name
      add_local_commit $main_branch_name 'conflicting_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_commit' 'conflicting_file' 'two'
      checkout_branch $feature_branch_name
      sha=`git log -n 1 | grep '^commit' | cut -d ' ' -f 2`
      git extract $new_branch_name $sha
    }

    it "aborts in the middle of the cherry-pick"
      expect_cherrypick_in_progress

    it "creates an abort script file"
      expect_file_exists "/tmp/git_extract_abort$temp_filename_suffix"

    function after {
      git cherry-pick --abort
    }

