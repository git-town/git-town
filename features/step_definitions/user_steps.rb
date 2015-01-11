Given(/^(I|my coworker) fetch(?:es)? updates$/) do |who|
  path = (who == 'I') ? local_repository_path : coworker_repository_path
  at_path path do
    run 'git fetch'
  end
end
