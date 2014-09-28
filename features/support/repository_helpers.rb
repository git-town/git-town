def repositiory_base
  '/tmp'
end

def remote_repository
  "#{repositiory_base}/git_town_specs_remote"
end

def coworker_repository
  "#{repositiory_base}/git_town_specs_coworker"
end

def local_repository
  "#{repositiory_base}/git_town_specs_local"
end

def create_repository path
  delete_repository path
  Dir.mkdir path
  run "git init --bare #{path}"
end

def clone_repository remote_path, path
  delete_repository path
  run "git clone #{remote_path} #{path}"
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
