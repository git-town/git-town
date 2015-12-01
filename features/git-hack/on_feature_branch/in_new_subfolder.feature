Feature: git hack: starting a new feature from a new subfolder on the main branch

  As a developer working on a feature branch that contains a subfolder that doesn't exist on the main branch
  I want to be able to extract my open changes into a new feature branch
  So that I can get them reviewed separately from the changes on this branch.


  This feature is untestable and unsupported by Git Town.
  GT performs all the correct commands here.
  But everything it does happens in a subshell, and doesn't affect the user shell.
  During this command, Git removes the folder that the user session is currently in,
  then later re-creates a new folder with the same name.
  The user session doesn't know that, so it is now in a folder that doesn't exist.

  When encountering this issue, simply cd into the Git root folder after git-hack is done.
