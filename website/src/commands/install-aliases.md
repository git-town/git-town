# git town aliases (add|remove)

The _aliases_ command adds or removes default global aliases. Global aliases
make Git Town commands feel like native Git commands. When enabled, you can run
`git hack` instead of `git town hack`. Please note that this can conflict with
other tools that also define Git aliases. This command does not overwrite
existing aliases. If you encounter issues, you can also set this manually for
individual commands:

```
git config --global alias.hack 'town hack'
```

### Variations

- when given `add`, creates the aliases
- when given `remove`, removes the aliases
