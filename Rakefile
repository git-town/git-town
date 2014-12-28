desc 'Run all linters and specs'
task default: %w(lint spec)


desc 'Run all formatters'
task format: %w(format:cucumber)

desc 'Run cucumber formatter'
task 'format:cucumber' do
  sh 'bundle exec cucumber_lint --fix'
end


desc 'Run all linters'
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


desc 'Run specs'
task :spec do
  sh 'bin/cuke'
end
