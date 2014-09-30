Given /^my repo has an upstream repo$/ do
  clone_repository remote_repository, upstream_repository, bare: true
  clone_repository upstream_repository, upstream_local_repository

  in_repository upstream_local_repository do
    run 'git checkout main'
  end

  run "git remote add upstream #{upstream_repository}"
end


Given /^my remote origin is "(.*)" on GitHub$/ do |repository|
  run "git remote set-url origin https://github.com/#{repository}.git"
end


Given /^my remote origin is a "rails\/rails" fork on GitHub through (.*)$/ do |protocol|
  url = case protocol
    when 'HTTPS' then github_rails_fork['clone_url']
    when 'SSH' then github_rails_fork['ssh_url']
    else raise "Unknown protocol: #{protocol}"
  end

  run "git remote set-url origin #{url}"
end



Then /^my remote upstream is "rails\/rails" on GitHub through (.*)$/ do |protocol|
  prefix = case protocol
    when 'HTTPS' then 'https://github.com/'
    when 'SSH' then 'git@github.com:'
    else raise "Unknown protocol: #{protocol}"
  end

  expect(remote_url('upstream')).to eql "#{prefix}rails/rails.git"
end
