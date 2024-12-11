#!/bin/bash

# Step 1: Switch to branch-2
git checkout branch-2

# Step 2: Rebase branch-2 onto main to remove branch-1's changes
git rebase --onto main branch-1 branch-2

# Step 3: Update branch-3 to continue from branch-2
git checkout branch-3
git rebase branch-2

# Verification
git log --graph --oneline --all
