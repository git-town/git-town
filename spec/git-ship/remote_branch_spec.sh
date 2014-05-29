require_remote_main_branch

describe 'git ship'

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
