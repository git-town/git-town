#!/usr/bin/env bash

# This script verifies that there are no files or folders that contain dashes.
# Git Town uses underscores in file paths.

files_with_dashes=$(find . -name '*-*' | grep -v node_modules | grep -v vendor | grep -v '.git' | grep -v website)
count=$(echo "$files_with_dashes" | wc -l)
if [ ! "$count" -eq 0 ]; then
  echo "Found $count files/folders containing dashes:"
  echo "$files_with_dashes"
fi
