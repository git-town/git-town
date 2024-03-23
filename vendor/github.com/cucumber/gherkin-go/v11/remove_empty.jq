# This jq script removes empty items. Theoretically Protobuf's JSON
# converter should do this, but it's inconsistent between language
# implementations.
#
# Original source:
# https://github.com/stedolan/jq/issues/104#issuecomment-338478029

def walk(f):
  . as $in
  | if type == "object" then
      reduce keys_unsorted[] as $key
        ( {}; . + { ($key):  ($in[$key] | walk(f)) } ) | f
  elif type == "array" then map( walk(f) ) | f
  else f
  end;
  
def remove_empty:
  . | walk(
    if type == "object" then
      with_entries(
        select(
          .value != null and
          .value != "" and
          .value != [] and
          .value != {}
        )
      )
    else .
    end
  );

remove_empty