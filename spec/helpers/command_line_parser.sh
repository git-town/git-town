# Tools to parse command line arguments given to the spec runner


# The name of the spec to run.
spec_filename=$1

if [ -z $spec_filename ]; then
  run_all_specs=true
else
  run_all_specs=false
fi

