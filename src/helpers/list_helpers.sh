#!/usr/bin/env bash

# Defines a list out of the given parameters
#
# Example: my_list=$(create_list 'hello world' two three)
#          returns a list with the elements "hello world, "two", and "three"
function create_list {
  for element in "$@"; do
    echo "$element"
  done
}


# Appends the given element to the list.
# If the appended element is a list, concatenates the two lists.
#
# Example: my_list=$(append_to_list "$my_list" five six)
#          returns a list containing the elements of $my_list
#          as well as "five" and "six" at the end
function append_to_list {
  while [ "$1" != "" ]; do
    echo "$1"
    shift
  done
}


# Returns a list that is the given list with the given elements removed.
#
# Example: filtered_list=$(remove_from_list "$my_list" two four)
#          returns a list containing the elements of $my_list, except the
#          elements "two" and "four"
function remove_from_list {
  local list=$1
  shift

  while [ "$1" != "" ]; do
    local element_to_remove=$1
    list=$(echo "$list" | grep -v "\b$element_to_remove\b")
    shift
  done
  echo "$list"
}


# Returns a new list that is the given list sorted.
# You can provide any argument that the UNIX command 'sort' understands.
#
# Example: sorted_list=$(sort_list "$my_list" -r)
#          returns a list that countains the elements of $my_list,
#          sorted in alphabetically reversed.
function sort_list {
  local list=$1
  shift

  echo "$list" | sort "$@"
}


# Returns the number of elements in the given list.
#
# Example: nr_of_elements=$(list_length "$my_list")
function list_length {
  local list=$1

  echo "$list" | wc -l | tr -d ' '
}


# Returns whether the given list contains the given element.
function list_contains {
  local list=$1
  local element=$2

  if [ "$(echo "$list" | grep -c "\b$element\b")" == "1" ]; then
    echo true
  else
    echo false
  fi
}


# Returns a list that is the given list truncated to the
# given number of elements, from the start.
# If the given index is negative, removes the given number of elements
# from the end of the list.
#
# Examples:
#  first_two_elements=$(first_n_elements_of_list "$my_list", 2)
#      returns the first two elements of $my_list
#  list_minus_last_two_elements=$(first_n_elements_of_list "$my_list", -2)
#      returns a list containing every element of $my_lis except
#      the last two elements
function first_n_elements_of_list {
  local list=$1
  local index=$2

  if [ "$index" -lt 0 ]; then
    local length && length=$(list_length "$list")
    index=$((length + index))
  fi
  echo "$list" | awk "FNR <= $index"
}


# Returns a list containing the last N numbers from the given list.
#
# Examples:
#  last_two_elements=$(last_n_elements_of_list "$my_list", 2)
#      returns the last two elements of $my_list
function last_n_elements_of_list {
  local list=$1
  local index=$2

  local length && length=$(list_length "$list")
  index=$((length - index))
  echo "$list" | awk "FNR > $index"
}


# Runs the given commands on each element of the given list,
# and returns the result as a new list.
#
# Example:
#   abbr=$(map_list "$my_list" "cut -c 1-2 | tr '[:lower:]' '[:upper:]'")
#   returns a list that contains the elements of $my_list,
#   truncated to the first 2 characters and uppercased.
function map_list {
  local list=$1
  shift
  eval "echo '$list' | xargs -L 1 echo | $*"
}



#############################
#
# Manual verification
#


function section {
  echo
  echo "$*"
}


section CREATING LISTS
numbers=$(create_list 'Kevin Goslar' two three)
echo "$numbers"
echo

section APPENDING ELEMENTS TO THE END
numbers=$(append_to_list "$numbers" four five)
echo "$numbers"

section REMOVING ELEMENTS
filtered=$(remove_from_list "$numbers" two four)
echo "$filtered"

section SORTING
sorted=$(sort_list "$numbers" -r)
echo "$sorted"

section PROCESSING
upper_cut=$(echo "$numbers" | xargs -L 1 echo | cut -c 1-2 | tr '[:lower:]' '[:upper:]')
echo "$upper_cut"
echo
upper_cut=$(map_list "$numbers" "cut -c 1-2 | tr '[:lower:]' '[:upper:]'")
echo "$upper_cut"


section COUNTING
count=$(list_length "$numbers")
echo "$count people in the house!"


section CONTAINS
echo "List contains two: $(list_contains "$numbers" two)"
echo "List contains zonk: $(list_contains "$numbers" zonk)"

section SLICING
first_two=$(first_n_elements_of_list "$numbers" 2)
echo "- the first two elements are: "
echo "$first_two"
echo

list_except_last_two=$(first_n_elements_of_list "$numbers" -3)
echo "- the list except the last three elements is: "
echo "$list_except_last_two"
echo

last_two=$(last_n_elements_of_list "$numbers" 2)
echo
echo "- the last two elements are: "
echo "$last_two"
