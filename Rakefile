desc 'Run all linters and tests'
task 'default' => %w(lint test)

desc 'Run all linters'
task 'lint' => %w(lint:bash lint:ruby)

desc 'Run the bash linter'
task 'lint:bash' do
  sh 'bin/lint'
end

desc 'Run the ruby linter'
task 'lint:ruby' do
  sh 'rubocop'
end

desc 'Run all tests'
task 'test' do
  sh 'cucumber'
end
