#!/usr/bin/env bash


# Prompts the user for the squash commit author. Sets the variable squash_commit_author
function get_squash_commit_author {
  local branch_name=$1
  local authors_with_commits="$(branch_authors "$branch_name")"
  local authors="$(echo "$authors_with_commits" | sed 's/^ *//' | cut -d ' ' -f 2-)"
  local number_of_authors="$(echo "$authors" | wc -l | tr -d ' ')"
  if [ "$number_of_authors" = 1 ]; then
    squash_commit_author="$(echo "$authors" | head -1)"
  else
    echo
    echo "Multiple people authored the '$branch_name' branch."
    echo "Please choose an author for the squash commit."
    echo

    local i=1
    echo "$authors_with_commits" | while read author_with_commits; do
      output_style_bold
      printf "%3s: " "$i"
      output_style_reset
      echo "$author_with_commits" |
        sed -E 's/^ *([0-9]+) *(.+)$/\2 (\1 commits)/' |
        sed -E 's/1 commits/1 commit/'
      i=$(( i + 1 ))
    done
    echo

    while [ -z "$squash_commit_author" ]; do
      echo -n "Enter user's number or a custom author (default: 1): "
      read input
      if [[ $input =~ ^[0-9]+$ ]]; then
        if [ "$input" -lt 1 ] || [ "$input" -gt "$number_of_authors" ]; then
          echo_inline_error "invalid number"
        else
          squash_commit_author="$(echo "$authors" | sed -n "${input}p")"
        fi
      elif [ -n "$input" ]; then
        squash_commit_author=$input
      else
        squash_commit_author="$(echo "$authors" | head -1)"
      fi
    done
    echo
  fi
}


# Returns the authors of the branch in the form "<commitCount> %an <%ae>"
function branch_authors {
  local branch_name=$1
  git shortlog -s -n -e "$MAIN_BRANCH_NAME..$branch_name" | tr '\t' '  '
}


# Returns the default author for new commits in the repository
function local_author {
  local name="$(git config user.name)"
  local email="$(git config user.email)"
  echo "$name <$email>"
}
