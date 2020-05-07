# Developing the Git Town source code

## adding a new Go library

- start using the new dependency in the code
- run `go mod vendor` to vendor it

## updating a single dependency

- `go get <path>`

## update all dependencies

<code textrun="verify-make-command">make update</code>
