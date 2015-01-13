desc 'Run linters and feature tests'
task default: %w(lint test)


desc 'Run formatters'
task format: %w(format:cucumber)

desc 'Run cucumber formatter'
task 'format:cucumber' do
  sh 'bundle exec cucumber_lint --fix'
end


desc 'Run linters'
task lint: %w(lint:bash lint:ruby lint:cucumber)

desc 'Run bash linter'
task 'lint:bash' do
  sh 'bin/lint'
end

desc 'Run ruby linter'
task 'lint:ruby' do
  sh 'bundle exec rubocop'
end

desc 'Run cucumber linter'
task 'lint:cucumber' do
  sh 'bundle exec cucumber_lint'
end


desc 'Run feature tests'
task :test do
  sh 'bin/cuke'
end
