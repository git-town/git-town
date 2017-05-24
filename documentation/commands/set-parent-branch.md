#### NAME

set-parent-branch - update the Git Town configuration for the parent branch of a feature branch


#### SYNOPSIS

```
git town set-parent-branch <child_branch_name> <parent_branch_name>
```


#### OPTIONS

```
<child_branch_name>
    The branch to update the parent for

<parent_branch_name>
    The new parent of <child_branch_name>
```


* *set-parent-branch*
> Update the parent branch of a feature branch
>
> ```bash
> # Set the parent branch of "feature-a" to "feature-b"
> git town set-parent-branch feature-a feature-b
> ```
