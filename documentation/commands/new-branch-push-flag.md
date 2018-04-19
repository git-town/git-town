<a textrun="command-heading">
#### New-branch-push-flag command
</a>

<a textrun="command-summary">
new-branch-push-flag - display or update the Git Town configuration for whether or not new branches are automatically pushed
</a>

<a textrun="command-description">
Branches created with hack / append / prepend will be pushed upon creation
if and only if "new-branch-push-flag" is true. The default value is false.
</a>

#### Usage

<a textrun="command-cli">
```
git town new-branch-push-flag [(true | false)]
git town new-branch-push-flag --global [(true | false)]
```
</a>


<a textrun="command-flags">
```
--global
    Display or update your global setting
```
</a>
