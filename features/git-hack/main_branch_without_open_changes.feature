Feature: hacking on main branch without open changes


  Background:
    Given the following commit exists in my repository
      | branch | location | message     | file name |
      | main   | remote   | main_commit | main_file |
    And I am on the main branch
    And I have an uncommitted file with name: "uncommitted" and content: "stuff"
    When I run `git hack feature`


  Scenario: on the main branch
    Then I end up on the "feature" branch
    And I still have an uncommitted file with name: "uncommitted" and content: "stuff"
    And the branch "feature" has not been pushed to the repository
    And I have the following commits
      | branch  | location         | message     | files     |
      | main    | local and remote | main_commit | main_file |
      | feature | local            | main_commit | main_file |
    And now I have the following committed files
      | branch  | files     |
      | main    | main_file |
      | feature | main_file |
