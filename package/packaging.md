# Git-Town Packaging and Release Process For Linux

This document contains instructions on how to update all the necessary packages once a new release is out.

The current linux (all distros) maintainer is [allonsy](https://github.com/allonsy). If an update occurs, please bug him

### Arch Linux
The package will need to be flagged out of date [here](https://aur.archlinux.org/packages/git-town/).
This will tell [allonsy](https://github.com/allonsy) to update the package (most likely minimal intervention necessary).
There is a button on the right hand side menu to flag the package.


### Debian Based Distributions
The package will need to be repackaged into a deb tarball.
There are two options: rebuild it from your own debian based machine or build it on the included docker image (recommended).

For either process you will need to update the version number for the package located in the following configuration files:
 * `package/debian/debian/changelog`


#### Process for your own machine:
1. Ensure that you have the following packages installed:
  * `build-essential`
  * `devscripts`
  * `debhelper`
2. Clone the repo (preferrably into a temporary directory to not mess with any existing branches)
3. cd into the `package/debian` directory
4. run the `./debian_build.sh` script
5. The deb file will be in the `package` directory
6. upload the deb file to the hosted repo (e.g. the releases page on github)

#### Process for the docker image:
1. First, fix the version variable in `package/debian/Dockerfile` to the correct value
  * Don't prepend a 'v', just put the version number
1. cd into `package/debian`
2. run `./docker_build.sh`
3. Enter the name of the folder you want to save the deb file to
3. Grab some coffee as the image builds and runs
4. The deb file should be stored in the given folder name gittown_version.deb


### Coming soon: rpm files for redhat/fedora based distros
