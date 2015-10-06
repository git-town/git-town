# A Flash of Bash

Git Town is written in Bash.
Bash is a surprisingly powerful, versatile, and elaborate scripting system
where every command-line application
becomes a keyword in your programming language.

With the right amount of structure,
Bash programs can be as easily written and read
as programs in many other languages.

Bash is the
[9th most popular language on GitHub](http://githut.info),
and is actively used in most open-source repositories.

Learning Bash not only allows you to contribute to Git Town,
it will make your life easier by giving you the ability to
create a plethora of useful little command-line helpers,
the UNIX way!

* [Shebang](#shebang)
* [If clauses](#if-clauses)
* [Functions](#functions)
* [Arrays](#arrays)
* [String Lists](#string-lists)
* [Converting](#converting)


## Shebang

Use the cross-platform version of
[shebang](https://en.wikipedia.org/wiki/Shebang_%28Unix%29):
```bash
#!/usr/bin/env bash
```


## If clauses

* structure

  ```bash
  if [ condition ]; then
    # code
  elif [ condition ]; then
    # code
  else
    # code
  fi
  ```

* negation

  ```bash
  if [ ! condition ]; then
    # runs if the condition is not true
  ```
* checking booleans

  ```bash
  if [ $success = true ]; then
    # runs if $success is true
  fi
  ```
* checking strings

  ```bash
  if [ $mood = "happy" ]; then
  ```
* checking for empty strings

  ```bash
  if [ -z $input ]; then
  ```

* checking whether a file exists

  ```bash
  if [ -e $filename ]; then
  ```


## Functions

* defining functions

  ```bash
  function add_user {
    local name=$1
    echo "adding user $name"
  }
  ```

  Parameter are named by assigning them
  to a local variable at the beginning of a function.

* calling functions

  ```bash
  add_user "git town"
  ```


## Arrays

* defining an array

  ```bash
  people=(Curly Larry Moe)
  ```
* adding elements to an array

  ```bash
  people+=(Mortimer)
  ```
* concatenating two arrays

  ```bash
  other_people=(Lucy)
  other_people+=("${people[@]}")
  ```
* cloning an array

  ```bash
  other_people=("${people[@]}")
  ```
* looping over the elements

  ```bash
  for person in "${people[@]}"; do
    echo "$person"
  done

  # looping with index
  for index in ${!people[*]}; do
    printf "%d: %s\n" $index ${people[$index]}
  done
  ```
* counting the elements

  ```bash
  count=${#people[*]}
  ```
* replacing an element

  ```bash
  people[1]="Mortimer"
  ```
* checking the existence of an array element

  ```bash
  function containsElement {
    local e
    for e in "${@:2}"; do [[ "$e" == "$1" ]] && echo true && return; done
    echo false
  }
  ```
* returning arrays from functions

  Functions can only return exit codes,
  no data structures.
  Use global variables to pass around arrays.


## String Lists

Lists are the universal exchange format between UNIX programs.
They are normal strings containing items separated by a newline character.
Unlike Bash arrays, lists allow to use the full power of the unix toolset for massaging them.

* defining

  ```bash
  people=$'curly\nlarry\nmoe'
  ```
* appending

  ```bash
  team=$(echo "$people" ; echo 'Judy')
  ```
* filtering / removing elements

  ```bash
  filtered=$(echo "$people" | grep -v curly)
  ```
* sorting

  ```bash
  sorted=$(echo "$people" | sort -r)
  ```
* processing

  ```bash
  upper=$(echo "$people" | tr '[:lower:]' '[:upper:]')
  ```
* counting

  ```bash
  count=$(echo "$people" | wc -l | tr -d ' ')
  ```


## Converting
* converting a string list into an array

  ```bash
  IFS=$'\n'
  people_array=($people)
  ```

* converting an array into a string list

  ```bash
  people_list=$( IFS=$'\n'; echo "${people_array[*]}" )
  ```
