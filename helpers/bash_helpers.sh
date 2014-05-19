# Helper methods for more convenient bash use

underline=`tput smul`
nounderline=`tput rmul`
bold=`tput bold`
normal=`tput sgr0`


# Determines whether a function with the given name exists
function determine_function_exists {
  if [ `type -t $1 | wc -l` == 1 ]; then
    function_exists=true
  else
    function_exists=false
  fi
}

