def repositiory_base
  '/tmp'
end

def remote_repository_path
  "#{repositiory_base}/git_town_specs_remote"
end

def coworker_repository_path
  "#{repositiory_base}/git_town_specs_coworker"
end

def local_repository_path
  "#{repositiory_base}/git_town_specs_local"
end

def upstream_remote_repository_path
  "#{repositiory_base}/git_town_specs_upstream_remote"
end

def upstream_local_repository_path
  "#{repositiory_base}/git_town_specs_upstream_local"
end

def create_repository path
  delete_repository path
  Dir.mkdir path
  run "git init --bare #{path}"
end

def clone_repository remote_path, path, options = {}
  delete_repository path
  run "git clone #{'--bare' if options[:bare]} #{remote_path} #{path}"
end

def at_path path
  cwd = Dir.pwd
  Dir.chdir path
  yield
  Dir.chdir cwd
end

def delete_repository path
  FileUtils.rm_r path, force: true
end
