require_main_branch

describe 'git ship'

  context 'simple local feature branch'

    function before {
      create_feature_branch $feature_branch_name
      add_local_commit $feature_branch_name 'feature_commit'
      git ship "feature_one"
    }


    it 'ends on the main branch'
      expect_current_branch_is $main_branch_name

    it 'squashes all commits into one on the main branch'
      expect_local_branch_has_commit $main_branch_name 'feature_one'

    it 'removes the feature branch'
      expect_no_local_branch_exists $feature_branch_name

    function after {
      git checkout $main_branch_name
      git branch -D $feature_branch_name
      git push origin :$feature_branch_name
    }

