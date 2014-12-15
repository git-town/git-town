#### NAME

git-town - view and change Git Town configuration and easily view help page


#### USAGE

View help screen:
```bash
git town
git town help
```

View the Git Town version:
```bash
git town version
```

View your current Git Town configuration:
```bash
git town config
```

View just your main branch configuration:
```bash
git town main-branch
```

Set your main branch to <branchname>:
```bash
git town main-branch <branchname>
```

View just your non-feature branches:
```bash
git town non-feature-branches
```

Add a new non-feature branch:
```bash
git town non-feature-branches --add <branchname>
```

Remove branch from non-feature branches:
```bash
git town non-feature-branches --remove <branchname>
```