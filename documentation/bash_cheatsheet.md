# A Flash of Bash

Git Town is written in Bash.
Bash is a surprisingly powerful, versatile, and elaborate scripting system
where every command-line application becomes a keyword in your programming language.

With the right amount of structure,
Bash programs can be as easily written and read
as programs in many other languages.

Bash is the
[9th most popular language on GitHub](http://githut.info),
and actively used in most open-source repositories.

Learning Bash not only allows you to contribute to Git Town,
it will make your life easier by giving you the ability to
create a plethora of useful little command-line helpers,
the UNIX way!


## Shebang

Use the cross-platform version of shebang:
```bash
#!/usr/bin/env bash
```


## Variables

* declaring variables

  ```bash
  name="Git Town"
  ```
* declaring local variables

  ```bash
  local name="Git Town"
  ```
* declaring constants:

  ```bash
  NAME="Git Town"
  ```


## Functions

* defining functions

  ```bash
  function add_user {
    local bar=$1
    echo "$bar"
  }
  ```

  Notice how the first parameter is given a name
  by assigning it
  to the local variable `bar`.

* calling functions

  ```bash
  add_user "git town"
  ```


## Arrays

* defining an array

  ```bash
  people=(Curly Larry Moe)
  ```

* replacing elements

  ```bash
  people[1]="Mortimer"
  ```

* counting the elements

  ```bash
  count=${#people[*]}
  ```

* cloning

  ```bash
  other_people=("${people[@]}")
  ```

* looping over the elements

  ```bash
  for person in "${people[@]}"; do
    echo $person
  done
  ```

  looping with index
  ```bash
  for index in ${!people[*]}; do
    printf "%4d: %s\n" $index ${people[$index]}
  done
  ```

* checking the existence of an array element

  ```bash
  function containsElement {
    local e
    for e in "${@:2}"; do [[ "$e" == "$1" ]] && echo true && return; done
    echo false
  }
  ```

* adding elements to an array

  ```bash
  people+=(Lucy)
  ```

## Lists

Lists are the universal exchange format between unix programs.
They are normal strings containing items separated by a newline character.
Unlike Bash arrays, lists allow to use the full power of the unix toolset for massaging them.

```bash
people=$'curly\nlarry\nmoe'
echo "$people"

echo FILTERING
filtered=$(echo "$people" | grep -v curly)
echo "$filtered"

echo SORTING
sorted=$(echo "$people" | sort -r)
echo "$sorted"

echo APPENDING
team=$(echo "$people" ; echo 'Judy')
echo "$team"

echo PROCESSING
upper=$(echo "$people" | tr '[:lower:]' '[:upper:]')
echo "$upper"

echo COUNTING
count=$(echo "$people" | wc -l | tr -d ' ')
echo "$count people in the house!"
```
