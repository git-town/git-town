#!/usr/bin/env sh
set -ex

# Template for submitting Git Town bug reproductions
#
# This script sets up a minimal Git repository in a directory named "test"
# and runs the Git Town command that exhibits the buggy behavior.
#
# To get you going, this script creates two branches and runs "git town sync".
# Please change this to create the conditions under which the bug happens.

echo "CREATE THE TEST REPO"
rm -rf test || true
git init test && cd test && git commit --allow-empty -m "initial"
git town status reset
git config set git-town.main-branch "main"
git config set git-town.new-branch-type "feature"
git config set git-town.sync-feature-strategy "rebase"

echo "CREATE BRANCH A"
git town hack branch-A
echo "content A" >file
git add file && git commit -m commit-A

echo "SYNC BRANCH A"
git town sync
