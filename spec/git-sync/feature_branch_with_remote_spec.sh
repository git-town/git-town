require_remote_main_branch

describe 'git-sync'

  context 'feature branch with remote branch'

    function before {
      git checkout -b $feature_branch_name
      git push -u origin $feature_branch_name
      add_local_commit $feature_branch_name 'local_feature_commit'
      add_remote_commit $feature_branch_name 'remote_feature_commit'
      add_local_commit $main_branch_name 'local_main_commit'
      add_remote_commit $main_branch_name 'remote_main_commit'
      checkout_feature_branch
      git sync
    }

    it 'remains on the feature branch'
      expect_current_branch_is $feature_branch_name

    it 'adds local updates from the main branch to the feature branch'
      expect_file_exists 'local_main_commit'

    it 'adds remote updates from the main branch to the feature branch'
      expect_file_exists 'remote_main_commit'

    it 'adds remote updates of the feature branch to the local feature branch'
      expect_file_exists 'remote_feature_commit'

    it 'pushes the feature branch to the repo'
      expect_synchronized_branch $feature_branch_name

