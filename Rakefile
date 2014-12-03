desc 'Run all linters and specs'
task 'default' => %w(lint spec)

desc 'Run all linters'
task 'lint' => %w(lint:bash lint:ruby)

desc 'Run bash linter'
task 'lint:bash' do
  sh 'bin/lint'
end

desc 'Run ruby linter'
task 'lint:ruby' do
  sh 'bundle exec rubocop'
end

desc 'Run specs'
task 'spec' do
  sh 'bundle exec parallel_cucumber features'
end
