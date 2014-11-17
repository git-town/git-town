# Git Town Release Notes


## 0.5
* **git sync-fork:** no longer automatically sets upstream configuration
  ([865030](https://github.com/Originate/git-town/commit/8650301a3ea40a989562a991960fa0d41b26f7f7))
* **linter for shell scripts**
  ([#149](https://github.com/Originate/git-town/issues/149))


## 0.4 (2014-11-13)

* **git kill:** completely removes a feature branch
  ([#87](https://github.com/Originate/git-town/issues/87),
   [edd7d8](https://github.com/Originate/git-town/commit/edd7d8180eb76717fd72e77d2c75edf8e3b6b6ca))
* **git sync:** pushes tags to the remote when running on the main branch
  ([#68](https://github.com/Originate/git-town/issues/68),
   [71b607](https://github.com/Originate/git-town/commit/71b607988c00e6dfc8f2598e9b964cc2ed4cfc39))
* **non-feature branches:** cannot be shipped and do not merge main when syncing
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


## 0.2 (2014-05-29)
* displays the duration of specs
* only pulls the main branch if it has a remote
* --abort options to abort failed Git Town operations
* --continue options to continue some Git Town operations after fixing the underlying issues
* can be installed through Homebrew
* colored test output
* display summary after tests
* exit with proper status codes
* better documentation


## 0.1 (2014-05-22)
* git hack, git sync, git extract, git ship
* basic test framework
* Travis CI integration
* self-hosting: uses Git Town for Git Town development

