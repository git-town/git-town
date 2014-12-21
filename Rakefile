desc 'Run all linters and specs'
task 'default' => %w(lint spec)

desc 'Run all linters'
task 'lint' => %w(lint:bash lint:ruby lint:cucumber)

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
  sh 'bundle exec cucumber_table_formatter -i'
end

desc 'Run specs'
task 'spec' do
  sh 'bin/cuke'
end
