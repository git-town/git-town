# Provides BDD-style test infrastructure


function describe {
  current_SUT=$1
  current_context=""
  current_spec_description=""
  function before_each { echo 'before_each not set'; }
  function after_each { echo 'after_each not set'; }
  function before { echo 'before not set'; }
  before_has_run=false
}


function context {
  current_context=$1
  current_spec_description=""
  function before_each { echo 'before_each not set'; }
  function after_each { echo 'after_each not set'; }
  before_has_run=false
}


function it {
  current_spec_description=$1
  if [ -z $current_context ]; then
    echo_header "$current_SUT $current_spec_description"
  else
    echo_header "$bold$current_SUT$normal $underline$current_context$nounderline $current_spec_description"
  fi
  if [ $before_has_run = false ]; then
    before
    before_has_run=true
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

