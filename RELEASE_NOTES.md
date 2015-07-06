# Git Town Release Notes

## 0.6.0 (2015-04-02)

* support for working without a remote repository for **git extract**, **git hack**, **git kill**, **git ship**, and **git sync**
  * implemented by our newest core committer @ricmatsui
* **git pr** renamed to **git pull-request**
  * set up an alias with `git config --global alias.pr pull-request`
* **git ship**
  * now accepts all `git commit` options
  * author with the most commits is automatically set as the author (when not the committer)
    ([#335](https://github.com/Originate/git-town/issues/335))
* **git pr/repo**
  * improved linux compatibility by trying `xdg-open` before `open`
* improved error messages when run outside a git repository
* improved setup wizard for initial configuration in a git repository
* added [contribution guide](/CONTRIBUTING.md)
* added [tutorial](/documentation/tutorial.md)


## 0.5.0 (2015-01-08)

* Manual installs need to update their `PATH` to point to the `src` folder within their clone of the repository
* **git extract:**
  * errors if branch exists remotely
    ([#236](https://github.com/Originate/git-town/issues/236))
  * removed restriction: need to be on a feature branch
    ([#269](https://github.com/Originate/git-town/issues/269))
  * added restriction: if no commits are provided, errors if the current branch does not have any have extractable commits (commits not in the main branch)
    ([#269](https://github.com/Originate/git-town/issues/269))
* **git hack:** errors if branch exists remotely
    ([#237](https://github.com/Originate/git-town/issues/237))
* **git kill:**
  * optional branch name
    ([#126](https://github.com/Originate/git-town/issues/126))
  * does not error if tracking branch has already been deleted
    ([#196](https://github.com/Originate/git-town/issues/196))
* **git pr:**
  * linux compatibility
    ([#232](https://github.com/Originate/git-town/issues/232))
  * compatible with more variants of specifying a Bitbucket or GitHub remote
    ([#271](https://github.com/Originate/git-town/issues/271))
  * compatible with respository names that contain ".git"
    ([#305](https://github.com/Originate/git-town/issues/305))
* **git repo:** view the repository homepage
  ([#140](https://github.com/Originate/git-town/issues/140))
* **git sync:**
  * `--all` option to sync all local branches
    ([#83](https://github.com/Originate/git-town/issues/83))
  * abort correctly after main branch updates and tracking branch conflicts
    ([#228](https://github.com/Originate/git-town/issues/228))
* **git town**: view and change Git Town configuration and easily view help page
  ([#98](https://github.com/Originate/git-town/issues/98))
* auto-completion for [Fish shell](http://fishshell.com)
  ([#177](https://github.com/Originate/git-town/issues/177))


## 0.4.1 (2014-12-02)

* **git pr:** create a new pull request
  ([#138](https://github.com/Originate/git-town/issues/138),
   [40d22e](https://github.com/Originate/git-town/commit/40d22eb1703ac96a58ac5052e70d20d7bdb9ac73))
* **git ship:**
  * empty commit message aborts the command
    ([#153](https://github.com/Originate/git-town/issues/153),
     [0bc84e](https://github.com/Originate/git-town/commit/0bc84ee626299896661fe1754cfa227630725bb9))
  * abort when there are no shippable changes
    ([#188](https://github.com/Originate/git-town/issues/188),
     [52fd94](https://github.com/Originate/git-town/commit/52fd94eca05bd3c2db5e7ac36121f08e56b9558b))
* **git sync:**
  * can now continue after just resolving conflicts (no need to commit or continue rebasing)
    ([#123](https://github.com/Originate/git-town/issues/123),
     [1a50ad](https://github.com/Originate/git-town/commit/1a50ad689a752f4eaed663e0ab22184621ee96a2))
  * restores deleted tracking branch
    ([#165](https://github.com/Originate/git-town/issues/165),
     [259464](https://github.com/Originate/git-town/commit/2594646ad853d83a6d697354d66755a374e42b8a))
* **git extract:** errors if branch already exists
  ([#128](https://github.com/Originate/git-town/issues/128),
   [75f498](https://github.com/Originate/git-town/commit/75f498771f19326f03bd1fd1bb70c9d9851b53f3))
* **git sync-fork:** no longer automatically sets upstream configuration
  ([865030](https://github.com/Originate/git-town/commit/8650301a3ea40a989562a991960fa0d41b26f7f7))
* remove needless checkouts for **git-ship**, **git-extract**, and **git-hack**
  ([#150](https://github.com/Originate/git-town/issues/150),
   [#155](https://github.com/Originate/git-town/issues/155),
   [8b385a](https://github.com/Originate/git-town/commit/8b385a745cf7ed28638e0a5c9c24440b7010354c),
   [35de43](https://github.com/Originate/git-town/commit/35de43156d9c6092840cd73456844b90acc36d8e))
* linters for shell scripts and ruby tests
  ([#149](https://github.com/Originate/git-town/issues/149),
   [076668](https://github.com/Originate/git-town/commit/07666825b5d60e15de274746fc3c26f72bd7aee2),
   [651c04](https://github.com/Originate/git-town/commit/651c0448309a376eee7d35659d8b06f709b113b5))
* rake tasks for development
  ([#170](https://github.com/Originate/git-town/issues/170),
   [ba74cf](https://github.com/Originate/git-town/commit/ba74cf30c8001941769dcd70410dbd18331f2fe9))


## 0.4.0 (2014-11-13)

* **git kill:** completely removes a feature branch
  ([#87](https://github.com/Originate/git-town/issues/87),
   [edd7d8](https://github.com/Originate/git-town/commit/edd7d8180eb76717fd72e77d2c75edf8e3b6b6ca))
* **git sync:** pushes tags to the remote when running on the main branch
  ([#68](https://github.com/Originate/git-town/issues/68),
   [71b607](https://github.com/Originate/git-town/commit/71b607988c00e6dfc8f2598e9b964cc2ed4cfc39))
* **perennial branches:** cannot be shipped and do not merge main when syncing
  ([#45](https://github.com/Originate/git-town/issues/45),
   [31dce1](https://github.com/Originate/git-town/commit/31dce1dfaf11e1e17f17e141a26cb38360ab731a))
* **git ship:**
  * merges main into the feature branch before squash merging
    ([#61](https://github.com/Originate/git-town/issues/61),
     [82d4d3](https://github.com/Originate/git-town/commit/82d4d3e745732cb397850a4c047826ba485e2bdb))
  * errors if the feature branch is not ahead of main
    ([#86](https://github.com/Originate/git-town/issues/86),
     [a0ace5](https://github.com/Originate/git-town/commit/a0ace5bb5e992c193df8abe4b0aca984c302c323))
  * git ship takes an optional branch name
    ([#95](https://github.com/Originate/git-town/issues/95),
     [cbf020](https://github.com/Originate/git-town/commit/cbf020fc3dd6d0ce49f8814a92f103e243f9cd2b))
* updated output to show each git command and its output, updated error messages
  ([8d8973](https://github.com/Originate/git-town/commit/8d8973aaa58394a123ceed2811271606f4e1aaa9),
   [60e1d8](https://github.com/Originate/git-town/commit/60e1d8299ebbb0e75bdae057e864d17e1f9a3ce7),
   [408e69](https://github.com/Originate/git-town/commit/408e699e5bdd3af524b2ea64669b81fea3bbe60b))
* skips unnecessary pushes
  ([0da896](https://github.com/Originate/git-town/commit/0da8968aef29f9ecb7326e0fafb5976f51789dca))
* **man pages**
  ([609e11](https://github.com/Originate/git-town/commit/609e11400818604328885df86c02ee4630410e12),
   [164f06](https://github.com/Originate/git-town/commit/164f06bc8bf00d9e99ce0416f408cf62959dc833),
   [27b257](https://github.com/Originate/git-town/commit/27b2573ca5ffa9ae7930f8b5999bbfdd72bd16d9))
* **git prune-branches**
  ([#48](https://github.com/Originate/git-town/issues/48),
   [7a922e](https://github.com/Originate/git-town/commit/7a922ecd9e03d20ed5a0c159022e601cebc80313))
* **Cucumber:** optional Fuubar output
  ([7c5402](https://github.com/Originate/git-town/commit/7c540284cf46bd49a7623566c1343285813524c6))


## 0.3 (2014-10-10)
* multi-user support for feature branches
  ([#35](https://github.com/Originate/git-town/issues/35),
   [ca0882](https://github.com/Originate/git-town/commit/ca08820c68457bddf6b8fff6c3ef3d430b905d9b))
* **git sync-fork**
  ([#22](https://github.com/Originate/git-town/issues/22),
   [1f1f9f](https://github.com/Originate/git-town/commit/1f1f9f98ffa7288d6a5982ec0c9e571695590fe1))
* stores configuration in the Git configuration instead of a dedicated file
  ([8b8695](https://github.com/Originate/git-town/commit/8b86953d7c7c719f28dbc7af6e86d02adaf2053e))
* only makes one fetch from the central repo per session
  ([#15](https://github.com/Originate/git-town/issues/15),
   [43400a](https://github.com/Originate/git-town/commit/43400a5b968a47eb55484f73e34026f66b1e939a))
* automatically prunes remote branches when fetching updates
  ([86100f](https://github.com/Originate/git-town/commit/86100f08866f19a0f4e80f470fe8dcc6996ddc2c))
* always cleans up abort and continue scripts after using one of them
  ([3be4c0](https://github.com/Originate/git-town/commit/3be4c06635a943f378287963ba30e4306fcd9802))
* simpler readme, dedicated RDD document
* **<a href="http://cukes.info" target="_blank">Cucumber</a>** feature specs (you need Ruby 2.x)
  ([c9d175](https://github.com/Originate/git-town/commit/c9d175fe2f28fbda3f662454f54ed80306ce2f46))
* much faster testing thanks to completely local test Git repos
  ([#25](https://github.com/Originate/git-town/issues/25),
   [c9d175](https://github.com/Originate/git-town/commit/c9d175fe2f28fbda3f662454f54ed80306ce2f46))


## 0.2.2 (2014-06-10)
* fixes "unary" error messages
* lots of output and documentation improvements


## 0.2.1 (2014-05-31)
* better terminal output
* Travis CI improvements
* better documentation


## 0.2.0 (2014-05-29)
* displays the duration of specs
* only pulls the main branch if it has a remote
* --abort options to abort failed Git Town operations
* --continue options to continue some Git Town operations after fixing the underlying issues
* can be installed through Homebrew
* colored test output
* display summary after tests
* exit with proper status codes
* better documentation


## 0.1.0 (2014-05-22)
* git hack, git sync, git extract, git ship
* basic test framework
* Travis CI integration
* self-hosting: uses Git Town for Git Town development

