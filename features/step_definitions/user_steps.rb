Given(/^my coworker fetches updates$/) do
  at_path coworker_repository_path do
    run 'git fetch'
  end
end
