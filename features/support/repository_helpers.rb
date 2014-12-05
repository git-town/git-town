def remote_repository_path
  "#{REPOSITORY_BASE}/git_town_specs_remote"
end

def coworker_repository_path
  "#{REPOSITORY_BASE}/git_town_specs_coworker"
end

def local_repository_path
  "#{REPOSITORY_BASE}/git_town_specs_local"
end

def upstream_remote_repository_path
  "#{REPOSITORY_BASE}/git_town_specs_upstream_remote"
end

def upstream_local_repository_path
  "#{REPOSITORY_BASE}/git_town_specs_upstream_local"
end

def create_repository path
  Dir.mkdir path
  run "git init --bare #{path}"
end

def clone_repository remote_path, path, bare: false
  run "git clone #{'--bare' if bare} #{remote_path} #{path}"
end

def at_path path
  cwd = Dir.pwd
  Dir.chdir path
  result = yield
  Dir.chdir cwd
  result
end
