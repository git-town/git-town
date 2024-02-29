# git set-parent

The _set-parent_ command changes the parent branch for the current branch. It
prompts the user for the new parent branch. Ideally you run [git sync](sync.md)
when done updating parent branches to pull the changes of the new parent
branches into their new child branches.

## Example

Let's say we have this branch hierarchy:

```
main
 |
 + feature-1
   |
   + feature-2
```

"feature-2" is a child branch of "feature-1". Let's make "feature-2" a child of
"main":

- run `git town set-parent`
- select `main` in the dialog

Now we have this branch hierarchy:

```
main
 |
 + feature-1
 |
 + feature-2
```
