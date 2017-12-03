# Git-Town Packaging and Release Process For Linux

This document contains instructions on how to update all the necessary packages once a new release is out.

The current linux (all distros) maintainer is [allonsy](https://github.com/allonsy). If an update occurs, please bug him

### Arch Linux

The package will need to be flagged out of date [here](https://aur.archlinux.org/packages/git-town/).
This will tell [allonsy](https://github.com/allonsy) to update the package (most likely minimal intervention necessary).
There is a button on the right hand side menu to flag the package.

### Debian Based Distributions

The process is now automated from travisCI.
When a new release is made from github, a new deb file is made and automatically uploaded to the releases page.

### Coming soon: rpm files for redhat/fedora based distros
