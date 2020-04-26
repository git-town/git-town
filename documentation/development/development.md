# Developing the Git Town source code

## update dependencies

<code textrun="verify-make-command">make update</code>

## adding a new Go library

- run `dep ensure --add <package name>`
- start using it in the code
- your pull request for the feature that requires the new library should contain
  the updated `Gopkg.*` files and vendor folder
