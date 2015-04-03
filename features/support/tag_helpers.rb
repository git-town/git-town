# Creates a tag with the given name in the local repo
def create_local_tag tag_name
  run "git tag -a #{tag_name} -m '#{tag_name}'"
end


# Creates a tag with the given name in the remote repo
def create_remote_tag tag_name
  in_secondary_repository do
    create_local_tag tag_name
    run 'git push --tags'
  end
end


# Returns the names of all unpushed tags.
def unpushed_tags
  output_of "git push --tags --dry-run 2>&1 | grep 'new tag' | awk '{print $4}'"
end
