#### NAME

new-branch-push-flag - display or update the Git Town configuration for whether or not new branches are automatically pushed

#### SYNOPSIS

```
git town new-branch-push-flag [(true | false)]
git town new-branch-push-flag --global [(true | false)]
```

#### DESCRIPTION

Branches created with hack / append / prepend will be pushed upon creation if and only if `new-branch-push-flag` is true. The default value is false.

#### OPTIONS

```
--global
    Display or update your global setting
```
