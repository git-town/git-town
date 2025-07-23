#!/usr/bin/env sh

# Template for submitting Git Town bug reproductions
#
# This script sets up a minimal Git repository in a directory named "test"
# and runs the Git Town command that exhibits the buggy behavior.
#
# To get you going, this script creates two branches and runs "git town sync".
# Please change this to create the conditions under which the bug happens.

set -x

rm -rf test || true

echo "CREATE THE TEST FOLDER"
git init test
cd test
git commit --allow-empty -m "initial"

echo "CONFIGURE GIT TOWN"
git config set git-town.main-branch "main"
git config set git-town.new-branch-type "feature"

echo "CREATE BRANCH A"
git town hack branch-A
echo "
line 1a
line 2a
" > file
git add file
git commit -m commit-A

echo "CREATE BRANCH B"
git town hack branch-B
echo "
line 1b
line 2b
" > file
git add file
git commit -m commit-B

echo "SHIP BRANCH A"
git checkout main
git merge --squash branch-A
git commit -m "branch A shipped"
git branch -d branch-A

echo "SYNC BRANCH B"
git checkout branch-B
git town sync
