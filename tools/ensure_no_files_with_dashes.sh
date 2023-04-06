#!/usr/bin/env bash

# This script verifies that there are no files or folders that contain dashes.
# Git Town uses underscores in file paths.

files_with_dashes=$(find . -name '*-*' | grep -v node_modules | grep -v vendor | grep -v '.git' | grep -v './website' | grep -v 'text-run.yml' | grep -v './.gherkin-*')
count=$(echo "$files_with_dashes" | wc -l)
if [ "$count" -gt 0 ]; then
  tput setaf 1
  echo
  echo "ERROR: Found $count files/folders containing dashes:"
  tput sgr0
  echo "$files_with_dashes"
  exit 1
fi
