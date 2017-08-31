Feature: git-new-pull-request when origin is on Bitbucket

  As a developer having finished a feature in a repository hosted on Bitbucket
  I want to be able to quickly create a pull request
  So that I have more time for coding the next feature instead of wasting it with process boilerplate.


  Scenario Outline: creating pull-requests
    Given I have a feature branch named "feature"
    And my remote origin is <ORIGIN>
    And I have "open" installed
    And I am on the "feature" branch
    When I run `git-town new-pull-request`
    Then I see a new pull request with the url "<URL>" in my browser

    Examples:
      | ORIGIN                                                            | URL                                                                                                                                                                        |
      | http://username@bitbucket.org/Originate/git-town.git              | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | http://username@bitbucket.org/Originate/git-town                  | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | https://username@bitbucket.org/Originate/git-town.git             | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | https://username@bitbucket.org/Originate/git-town                 | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | git@bitbucket.org/Originate/git-town.git                          | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | git@bitbucket.org/Originate/git-town                              | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | git@bitbucket-as-account2:Originate/git-town.git                  | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | http://username@bitbucket.org/Originate/originate.github.com.git  | https://bitbucket.org/Originate/originate.github.com/pull-request/new?dest=Originate%2Foriginate.github.com%3A%3Amain&source=Originate%2Foriginate.github.com%.*%3Afeature |
      | http://username@bitbucket.org/Originate/originate.github.com      | https://bitbucket.org/Originate/originate.github.com/pull-request/new?dest=Originate%2Foriginate.github.com%3A%3Amain&source=Originate%2Foriginate.github.com%.*%3Afeature |
      | https://username@bitbucket.org/Originate/originate.github.com.git | https://bitbucket.org/Originate/originate.github.com/pull-request/new?dest=Originate%2Foriginate.github.com%3A%3Amain&source=Originate%2Foriginate.github.com%.*%3Afeature |
      | https://username@bitbucket.org/Originate/originate.github.com     | https://bitbucket.org/Originate/originate.github.com/pull-request/new?dest=Originate%2Foriginate.github.com%3A%3Amain&source=Originate%2Foriginate.github.com%.*%3Afeature |
      | git@bitbucket.org/Originate/originate.github.com.git              | https://bitbucket.org/Originate/originate.github.com/pull-request/new?dest=Originate%2Foriginate.github.com%3A%3Amain&source=Originate%2Foriginate.github.com%.*%3Afeature |
      | git@bitbucket.org/Originate/originate.github.com                  | https://bitbucket.org/Originate/originate.github.com/pull-request/new?dest=Originate%2Foriginate.github.com%3A%3Amain&source=Originate%2Foriginate.github.com%.*%3Afeature |
      | ssh://git@bitbucket.org/Originate/git-town.git                    | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |
      | ssh://git@bitbucket.org/Originate/git-town                        | https://bitbucket.org/Originate/git-town/pull-request/new?dest=Originate%2Fgit-town%3A%3Amain&source=Originate%2Fgit-town%.*%3Afeature                                     |

