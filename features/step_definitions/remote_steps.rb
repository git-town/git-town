# frozen_string_literal: true

Given(/^my repo has an upstream repo$/) do
  clone_repository :origin, :upstream, bare: true
  clone_repository :upstream, :upstream_developer
  run "git remote add upstream #{repository_path :upstream}"
end


Given(/^my repo does not have a remote origin$/) do
  run 'git remote rm origin'
end


Given(/^my repo's remote origin is (.+?)$/) do |origin|
  run "git config git-town.testing.remote-url #{origin}"
end
