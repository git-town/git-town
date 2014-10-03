Given /^my coworker Charlie works on the same feature branch$/ do
  at_path coworker_repository_path do
    run "git pull ; git checkout feature"
  end
end

