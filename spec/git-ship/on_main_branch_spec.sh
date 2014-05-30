require_main_branch

describe 'git ship'

  context 'on the main branch'

    function before {
      git ship
    }


    it 'remains on the main branch'
      expect_current_branch_is $main_branch_name

    it 'does nothing'
      expect_local_branch_count 2

