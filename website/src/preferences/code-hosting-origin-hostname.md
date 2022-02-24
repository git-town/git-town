# code-hosting-origin-hostname

```
git-town.code-hosting-origin-hostname=<hostname>
```

When using SSH identities, you can use this configuration setting to define the
hostname of your source code repository by running:

```
git config [--global] git-town.code-hosting-origin-hostname <hostname>
```

`<hostname>` should match the hostname in your ssh config file. The optional
`--global` flag applies this setting to all Git repositories on your local
machine. When not present, the setting applies to the current repo.
