Given(/^my repo has an upstream repo$/) do
  clone_repository :origin, :upstream, bare: true
  clone_repository :upstream, :upstream_developer
  run "git remote add upstream #{repository_path :upstream}"
end


Given(/^my remote origin is (.+?)$/) do |origin|
  run "git remote set-url origin #{origin}"
end
