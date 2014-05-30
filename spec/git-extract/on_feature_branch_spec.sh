describe "git extract"

  context "on a feature branch"

    function before {
      require_remote_main_branch
      new_branch_name="extracted_branch_$$_$current_time"
      create_feature_branch $feature_branch_name
      add_local_commit $feature_branch_name 'feature_commit'
      add_remote_commit $main_branch_name 'remote_main_commit'
      checkout_branch $feature_branch_name
      sha=`git log -n 1 | grep '^commit' | cut -d ' ' -f 2`
      git extract $new_branch_name $sha
    }


    it 'ends up on the new created branch'
      expect_current_branch_is $new_branch_name

    it 'pushes the new branch to the repo'
      expect_synchronized_branch $new_branch_name

    it 'updates the main branch'
      expect_local_branch_has_commit $main_branch_name 'remote_main_commit'

    it 'bases the created feature branch off the updated main branch'
      expect_local_branch_has_commit $new_branch_name 'remote_main_commit'
