describe 'git-sync'

  context 'on the main branch'

    function before {
      require_remote_main_branch
      add_remote_commit $main_branch_name 'remote_main_commit'
      add_local_commit $main_branch_name 'local_main_commit'
      git checkout $main_branch_name
      git sync

    }

    it 'pulls updates from the remote branch'
      expect_file_exists 'remote_main_commit'

    it 'pushes local updates to the repo'
      expect_synchronized_branch $main_branch_name

