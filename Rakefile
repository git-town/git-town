# frozen_string_literal: true

require 'active_support/all'
require 'rubocop/rake_task'


desc 'Run linters and feature tests'
task default: %w[lint unit-test feature-test]


# Formatters
desc 'Run formatters'
task format: %w[rubocop:auto_correct format:cucumber]

desc 'Run cucumber formatter'
task 'format:cucumber' do
  sh 'bundle exec cucumber_lint --fix'
end


# Linters
desc 'Run linters'
task lint: %w[lint:go lint:ruby lint:cucumber]

desc 'Run Go linter'
task 'lint:go' do
  sh 'bin/lint_go'
end

desc 'Run Ruby linter'
task 'lint:ruby' => [:rubocop]

desc 'Run Cucumber linter'
task 'lint:cucumber' do
  sh 'bundle exec cucumber_lint'
end


# Unit tests
desc 'Run unit tests'
task 'unit-test' do
  sh 'bin/tests'
end

# Feature tests
desc 'Run feature tests'
task 'feature-test' do
  sh 'bin/cuke'
end

def run command
  sh command
  puts
end

desc 'Deploys a new version of the website'
task 'deploy' do
  run 'git checkout gh-pages'
  run 'git pull'
  run 'git checkout master'
  run 'git pull --rebase'
  run 'harp compile website/ _www'
  run 'git checkout gh-pages'
  run 'cp -r _www/* .'
  run 'rm -rf _www'
  run 'git add -A'
  print 'Description of this change: '
  desc = STDIN.gets.strip
  return if desc.blank?
  run "git commit -m '#{desc}'"
  run 'git push'
  run 'git checkout master'
end


RuboCop::RakeTask.new
