#### NAME

git-town - view and change Git Town configuration and easily view help page


#### SYNOPSIS

```
git town
git town config
git town help
git town main-branch [<branchname>]
git town non-feature-branches [(--add | --remove) <branchname>]
git town version
```

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

View just your main branch or non-feature branch configuration:
```bash
git town main-branch
git town non-feature-branches
```

Set your main branch to <branchname>:
```bash
git town main-branch <branchname>
```

Add/remove a non-feature branch:
```bash
git town non-feature-branches --add    <branchname>
git town non-feature-branches --remove <branchname>
```