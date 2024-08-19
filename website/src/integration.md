# Integration

This page describes how to integrate Git Town into other applications.

## [Lazygit](https://github.com/jesseduffield/lazygit)

Example lazygit configuration file to integrate Git Town:

```yml
customCommands:
  - key: 'Y'
    context: 'global'
    description: 'Git-Town sYnc'
    command: 'git-town sync --all'
    stream: true
    loadingText: 'Syncing'
  - key: 'U'
    context: 'global'
    description: 'Git-Town Undo (undo the last git-town command)'
    command: 'git-town undo'
    prompts:
    - type: 'confirm'
      title: 'Undo Last Command'
      body: 'Are you sure you want to Undo the last git-town command?'
    stream: true
    loadingText: 'Undoing Git-Town Command'
  - key: '!'
    context: 'global'
    description: 'Git-Town Repo (opens the repo link)'
    command: 'git-town repo'
    stream: true
    loadingText: 'Opening Repo Link'
  - key: 'a'
    context: 'localBranches'
    description: "Git-Town Append"
    prompts:
      - type: 'input'
        title: "Enter name of new child branch. Branches off of '{{.CheckedOutBranch.Name}}'"
        key: 'BranchName'
    command: 'git-town append {{.Form.BranchName}}'
    stream: true
    loadingText: 'Appending'
  - key: 'h'
    context: 'localBranches'
    description: 'Git-Town Hack (creates a new branch)'
    prompts:
      - type: 'input'
        title: "Enter name of new branch. Branches off of 'Main'"
        key: 'BranchName'
    command: 'git-town hack {{.Form.BranchName}}'
    stream: true
    loadingText: 'Hacking'
  - key: 'K'
    context: 'localBranches'
    description: 'Git-Town Kill (deletes the current feature branch and sYnc)'
    command: 'git-town kill'
    prompts:
    - type: 'confirm'
      title: 'Delete current feature branch'
      body: 'Are you sure you want to delete the current feature branch?'
    stream: true
    loadingText: 'Killing Feature Branch'
  - key: 'p'
    context: 'localBranches'
    description: 'Git-Town Propose (creates a pull request)'
    command: 'git-town propose'
    stream: true
    loadingText: 'Creating pull request'
  - key: 'P'
    context: 'localBranches'
    description: "Git-Town Prepend (creates a branch between the curent branch and its parent)"
    prompts:
      - type: 'input'
        title: "Enter name of the for child branch between '{{.CheckedOutBranch.Name}}' and its parent"
        key: 'BranchName'
    command: 'git-town prepend {{.Form.BranchName}}'
    stream: true
    loadingText: 'Prepending'
  - key: 'S'
    context: 'localBranches'
    description: 'Git-Town Skip (skip branch with merge conflicts when syncing)'
    command: 'git-town skip'
    stream: true
    loadingText: 'Skiping'
  - key: 'G'
    context: 'files'
    description: 'Git-Town GO aka:continue (continue after resolving merge conflicts)'
    command: 'git-town continue'
    stream: true
    loadingText: 'Continuing'
```
