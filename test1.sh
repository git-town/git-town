#!/bin/bash

# Initialize a new Git repository
cd ~
mkdir demo-repo
cd demo-repo
git init

# Create an initial commit on the main branch
echo "Initial content" > file.txt
git add file.txt
git commit -m "Initial commit"

# Create branch-1 with some changes
git checkout -b branch-1
echo "Change from branch-1" >> file1.txt
git add file1.txt
git commit -m "commit 1"

# Create branch-2 based on branch-1 with additional changes
git checkout -b branch-2
echo "Change from branch-2" >> file2.txt
git add file2.txt
git commit -m "commit 2"

# Create branch-3 based on branch-2 with additional changes
git checkout -b branch-3
echo "Change from branch-3" >> file3.txt
git add file3.txt
git commit -m "commit 3"

echo "Repository created in the 'before' state"
