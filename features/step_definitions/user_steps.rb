Given(/^my coworker fetches updates$/) do
  in_repository :coworker do
    run 'git fetch'
  end
end
