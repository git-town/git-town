Given /^the following commits? exists?$/ do |commits_data|

  # Save the current branch in order to restore it later
  current_branch = run_this("git rev-parse --abbrev-ref HEAD")[:out]

  commits_data.hashes.each do |commit_data|

    # Gather all the given options and augment with default values
    options = {
      file_name: commit_data.delete('file name') { 'default file name' },
      file_content: commit_data.delete('file content') { 'default file content' },
      commit_message: commit_data.delete('message') { 'default commit message' },
      commit_location: commit_data.delete('location'){%i[local remote]},
      branch: commit_data.delete('branch')
    }
    if options[:commit_location].is_a? String
      options[:commit_location] = [options[:commit_location].to_sym]
    end

    # Make sure we understood all commit data
    if commit_data != {}
      raise "Unused commit specifiers: #{commit_data}"
    end

    # Check out the respective branch
    run_this "git checkout #{options[:branch]}", allow_failures: true

    # Create commits
    if options[:commit_location].delete :local
      create_local_commit options
    end
    if options[:commit_location].delete :remote
      IO.write options[:file_name], options[:file_content]
      run_this "git add #{options[:file_name]} ; git commit -m '#{options[:commit_message]}' ; git push ; git reset --hard HEAD^"
    end

    if options[:commit_location] != []
      raise "Unused commit location: #{options[:commit_location]}"
    end
  end

  # Go back to the branch that was checked out initially
  run_this "git checkout #{current_branch}", allow_failures: true
end




Then /^my branch and its remote still have (\d+) and (\d+) different commits$/ do |local_count, remote_count|
  matches = /have (\d+) and (\d+) different commit each/.match(run_this("git status")[:out])
  expect(matches[1]).to eql local_count
  expect(matches[2]).to eql remote_count
end


Then /^(now )?I (still )?have the following commits$/ do |_, _, commits_data|
  expected_commits = commits_data.hashes
                                 .each do |commit_data|
                                    symbolize_keys_deep! commit_data
                                    commit_data[:files] = commit_data[:files].split(',')
                                                                             .map(&:strip)
                                  end
  expect(actual_commits).to match_array expected_commits
end


Then /^there are no commits$/ do
  expect(actual_commits).to eql []
end
