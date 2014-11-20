Given(/^my repo has an upstream repo$/) do
  clone_repository remote_repository_path, upstream_remote_repository_path, bare: true
  clone_repository upstream_remote_repository_path, upstream_local_repository_path

  at_path upstream_local_repository_path do
    run 'git checkout main'
  end

  run "git remote add upstream #{upstream_remote_repository_path}"
end
