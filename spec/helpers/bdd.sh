# Provides BDD-style test infrastructure


function assert {
  actual_value=$1
  expected_value=$2
  name=$3

  if [ ! "$actual_value" = "$expected_value" ]; then
    if [ -z $name ]; then
      echo_failure "Expected '$expected_value', but got '$actual_value'"
    else
      echo_failure "Expected '$name' to equal '$expected_value', but it was '$actual_value'"
    fi
  else
    if [ -z $name ]; then
      echo_success "Found '$actual_value' as expected"
    else
      echo_success "'$name' is as expected '$actual_value'"
    fi
  fi
}


function describe {
  current_SUT=$1
  current_context=""
  current_spec_description=""
  function before_each { echo 'before_each not set'; }
  function after_each { echo 'after_each not set'; }
}


function context {
  current_context=$1
  current_spec_description=""
  function before_each { echo 'before_each not set'; }
  function after_each { echo 'after_each not set'; }
}


function it {
  current_spec_description=$1
  if [ -z $current_context ]; then
    echo_header "$current_SUT $current_spec_description"
  else
    echo_header "$current_SUT $current_context $current_spec_description"
  fi
  before_each
}


# Finishes an "it" block.
#
# Calling this method is only required if there are "after_each" blocks that
# should be executed.
function ti {
  after_each
}

