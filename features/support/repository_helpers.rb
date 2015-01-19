# Execute the block at the given path
def at_path path
  cwd = Dir.pwd
  Dir.chdir path
  result = yield
  Dir.chdir cwd
  result
end


# Create a repository for the given identifier
def create_repository identifier, memoized: false
  path = repository_path identifier, memoized: memoized
  Dir.mkdir path
  run "git init --bare #{path}"
end


# Clone the repository signified by parent_idenifier into child_identifier
def clone_repository parent_identifier, child_identifier, bare: false, memoized: false
  parent_path = repository_path parent_identifier, memoized: memoized
  child_path = repository_path child_identifier, memoized: memoized

  run "git clone #{'--bare' if bare} #{parent_path} #{child_path}"

  in_repository child_identifier, memoized: memoized do
    user = child_identifier.to_s.sub('_secondary', '')
    configure_git user
  end
end


# Move into the repository with the given repository
def go_to_repository identifier
  Dir.chdir repository_path identifier
end


# Execute the block in the repository for the given identifier
def in_repository identifier, parent: :origin, memoized: false, &block
  path = repository_path identifier, memoized: memoized
  clone_repository parent, identifier, memoized: memoized unless File.directory? path
  at_path path, &block
end


# Execute the block in a secondary repository of the current user
def in_secondary_repository &block
  current_idenitifer = Pathname.new(Dir.pwd).basename
  in_repository "#{current_idenitifer}_secondary", &block
end


# Returns the repository path for the given identifier
def repository_path identifier, memoized: false
  "#{memoized ? MEMOIZED_REPOSITORY_BASE : REPOSITORY_BASE}/#{identifier}"
end
