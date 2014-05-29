require_main_branch

describe 'git ship'

  context 'on the main branch'

    function before {
      git ship
    }

    it 'does nothing'
      expect_local_branch_count 2

