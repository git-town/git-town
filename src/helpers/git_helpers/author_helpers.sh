#!/usr/bin/env bash


# Prompts the user for the squash commit author. Sets the variable squash_commit_author
function get_squash_commit_author {
  local branch_name=$1
  local authors="$(branch_authors "$branch_name")"
  local number_of_authors="$(echo "$authors" | wc -l | tr -d ' ')"
  if [ "$number_of_authors" = 1 ]; then
    squash_commit_author="$(remove_author_commits "$(echo "$authors" | head -1)")"
  else
    echo
    echo "Multiple people authored the '$branch_name' branch."
    echo "Please choose an author for the squash commit."
    echo

    local i=1
    echo "$authors" | while read author; do
      output_style_bold
      printf "%3s: " "$i"
      output_style_reset
      echo "$author"
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
          squash_commit_author="$(remove_author_commits "$(echo "$authors" | sed -n "${input}p")")"
        fi
      elif [ -n "$input" ]; then
        squash_commit_author=$input
      else
        squash_commit_author="$(remove_author_commits "$(echo "$authors" | head -1)")"
      fi
    done
    echo
  fi
}


# Returns the authors of the branch in the form <commits> <author name and email>
function branch_authors {
  local branch_name=$1
  git log "$MAIN_BRANCH_NAME..$branch_name" --format='%an <%ae>' | # Authors of commits only in $branch_name
    uniq -c  | # get count for each author
    sort -rn | # reverse numeric sort
    sed -E 's/^ *([0-9]+) (.+)$/\2 (\1 commits)/g' | # transform to "<author> (# commits)"
    sed -E 's/1 commits/1 commit/g'
}


# Returns the default author for new commits in the repository
function local_author {
  local name="$(git config user.name)"
  local email="$(git config user.email)"
  echo "$name <$email>"
}


# Removes the the number of commits from an author string
function remove_author_commits {
  local author=$1
  echo "$author" | sed -E 's/ \([0-9]* commits?\)$//'
}
