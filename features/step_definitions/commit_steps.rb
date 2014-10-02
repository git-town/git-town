Given /^the following commits? exists? in my repository$/ do |commits_table|
  at_path local_repository_path do
    create_commits commits_table
  end
end


Given /^the following commits exist in Charly's repository$/ do |commits_table|
  at_path coworker_repository_path do
    create_commits commits_table
  end
end




Then /^my branch and its remote still have (\d+) and (\d+) different commits$/ do |local_count, remote_count|
  matches = /have (\d+) and (\d+) different commit each/.match(run("git status")[:out])
  expect(matches[1]).to eql local_count
  expect(matches[2]).to eql remote_count
end


Then /^(?:now )?(?:(?:I (?:still )?(?:have|see))) the following commits$/ do |commits_table|
  expected_commits = commits_table.hashes
                                  .each do |commit_data|
                                    symbolize_keys_deep! commit_data
                                    commit_data[:files] = commit_data[:files].split(',')
                                                                             .map(&:strip)
                                    commit_data[:location] = Kappamaki.from_sentence commit_data[:location]
                                  end
  expected_commits.map! do |commit_data|
    locations = commit_data.delete :location
    locations.map do |location|
      result = commit_data.clone
      result[:location] = location
      if location == 'remote' && /^[^\/]+$/.match(result[:branch])
        result[:branch] = "remotes/origin/#{result[:branch]}"
      end
      result
    end
  end.flatten!

  at_path local_repository_path do
    expect(commits_in_repo).to match_commits expected_commits
  end
end


Then /^(?:now )?Charly(?: still)? sees the following commits$/ do |commits_table|
  expected_commits = commits_table.hashes
                                  .each do |commit_data|
                                    symbolize_keys_deep! commit_data
                                    commit_data[:files] = commit_data[:files].split(',')
                                                                             .map(&:strip)
                                    commit_data[:location] = Kappamaki.from_sentence commit_data[:location]
                                  end
  expected_commits.map! do |commit_data|
    locations = commit_data.delete :location
    locations.map do |location|
      result = commit_data.clone
      result[:location] = location
      if location == 'remote' && /^[^\/]+$/.match(result[:branch])
        result[:branch] = "remotes/origin/#{result[:branch]}"
      end
      result
    end
  end.flatten!

  at_path coworker_repository_path do
    expect(commits_in_repo).to match_commits expected_commits
  end
end


Then /^now the following commits exist$/ do |commits_data|
  expected_commits = commits_data.hashes
                                 .each do |commit_data|
                                    symbolize_keys_deep! commit_data
                                    commit_data[:files] = commit_data[:files].split(',')
                                                                             .map(&:strip)
                                  end
  all_existing_commits = []
  at_path local_repository_path do
    all_existing_commits.concat actual_commits.each {|commit| commit[:location] = 'local'}
    run "git checkout origin/feature"
    all_existing_commits.concat actual_commits.each {|commit| commit[:location] = 'remote'}
  end
  at_path coworker_repository_path do
    all_existing_commits.concat actual_commits.each {|commit| commit[:location] = 'Charly'}
  end
  expect(all_existing_commits).to match_array expected_commits
end


Then /^there are no commits$/ do
  expect(commits_in_repo).to eql []
end
