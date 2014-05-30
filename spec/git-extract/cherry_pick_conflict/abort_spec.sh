require_remote_main_branch

describe "git-extract with conflicts while cherry-picking"

  context "aborting the extract"

    function before {
      new_branch_name="extracted_branch_$$_$current_time"
      create_feature_branch $feature_branch_name
      add_local_commit $main_branch_name 'conflicting_commit' 'conflicting_file' 'one'
      add_local_commit $feature_branch_name 'conflicting_commit' 'conflicting_file' 'two'
      checkout_branch $feature_branch_name
      sha=`git log -n 1 | grep '^commit' | cut -d ' ' -f 2`
      git extract $new_branch_name $sha
      git extract --abort
    }

    it "returns to the feature branch"
      expect_current_branch_is $feature_branch_name

    it "removes the extracted branch"
      expect_no_local_branch_exists $new_branch_name

    it "stops the cherry-pick"
      expect_no_cherrypick_in_progress
