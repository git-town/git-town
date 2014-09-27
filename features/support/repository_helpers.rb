def repositiory_base_path
  '/tmp'
end

def remote_repository_path
  "#{repositiory_base_path}/git_town_specs_remote"
end

def coworker_repository_path
  "#{repositiory_base_path}/git_town_specs_coworker"
end

def local_repository_path
  "#{repositiory_base_path}/git_town_specs_local"
end

def create_repository path, &block
  delete_repository path
  Dir.mkdir path
  run "git init --bare #{path}"
  in_repository path, &block
end

def clone_repository remote_path, path, &block
  delete_repository path
  run "git clone #{remote_path} #{path}"
  in_repository path, &block
end

def in_repository path, &block
  cwd = Dir.pwd
  Dir.chdir path
  yield if block_given?
  Dir.chdir cwd
end

def delete_repository path
  FileUtils.rm_r path, force: true
end
