require_remote_main_branch

describe 'git ship'

  context 'on the main branch'

    function before {
      git ship
    }

    it 'does nothing'
      expect_local_branch_count 2


  context 'simple local feature branch'

    function before {
      git checkout -b feature1
      add_local_commit "feature1" 'feature_commit'
      git ship "feature_one"
    }


    it 'ends on the main branch'
      expect_current_branch_is $main_branch_name

    it 'squashes all commits into one on the main branch'
      expect_local_branch_has_commit $main_branch_name 'feature_one'

    it 'removes the feature branch'
      expect_no_local_branch_exists 'feature1'

    function after {
      git checkout $main_branch_name
      git branch -D feature1
      git push origin :feature1
    }


  context 'with non-pulled updates on the remote branch'

    function before {
      git checkout -b feature1
      git push -u origin feature1
      add_remote_commit 'feature1' 'remote_commit'
      git ship 'feature_one'
    }

    it 'includes the remote commits into the merge'
      expect_file_content 'remote_commit' 'remote_commit'

    it 'removes the remote branch'
      expect_no_remote_branch_exists 'feature1'

    it 'removes the feature branch'
      expect_no_local_branch_exists 'feature1'
