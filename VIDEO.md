# Creating animations of the development work

You can use [gource](https://gource.io) to create an animation of the commit
activity.

For a nicer video, the video should not include the "vendor" folder since we
don't really "work" on the files in there. To do that:

- create a copy of the repo

  ```bash
  cp -r git-town/ git-town-video/
  ```
- enter the new folder

  ```bash
  cd git-town-video
  ```
- remove the "vendor" folder from all commits

  ```bash
  git filter-branch --tree-filter 'rm -rf vendor' HEAD
  ```
- run gource

  ```bash
  gource --load-config .gource.conf
  ```
