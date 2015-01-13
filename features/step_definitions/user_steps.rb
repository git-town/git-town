Given(/^(I|my coworker) fetch(?:es)? updates$/) do |who|
  user = (who == 'I') ? :developer : :coworker
  in_repository user do
    run 'git fetch'
  end
end
