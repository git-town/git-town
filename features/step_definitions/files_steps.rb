Given /^I have an uncommitted file with name: "(.*?)" and content: "(.*?)"$/ do |name, content|
  IO.write name, content
end


Then /^a the main branch name is configured as "(.*?)"$/ do |main_branch_name|
  expect(IO.read '.gittownrc').to eql "#{main_branch_name}\n"
end


Then /^file "(.*?)" has a merge conflict$/ do |file|
  run("git status | grep 'both added.*#{file}' | wc -l")[:out] == '1'
end


Then /^there are no merge conflicts anymore$/ do
  run("git status | grep 'both added' | wc -l")[:out] == '0'
end


Then /^(now I|I still) have the following committed files$/ do |_, files_data|

  # Get all files in all branches
  actual_files = []
  existing_local_branches.each do |branch|
    run "git checkout #{branch}"
    existing_files.each do |file|
      if file != "uncommitted"
        actual_files << {branch: branch, name: file, content: IO.read(file)}
      end
    end
  end

  # Get all expected files
  expected_files = files_data.hashes.each do |expected_file|
    symbolize_keys_deep! expected_file
    if expected_file[:branch] == 'main'
      expected_file[:branch] = 'main'
    end
  end

  # Remove the keys that are not used in the expected data
  used_keys = expected_files[0].keys
  actual_files.each do |actual_file|
    actual_file.keys.each do |key|
      unless used_keys.include? key
        actual_file.delete key
      end
    end
  end

  expect(actual_files).to match_array expected_files
end




Then /^I don't have an uncommitted file with name: "(.*?)"$/ do |file_name|
  actual_files = run("git status --porcelain | awk '{print $2}'")[:out].split("\n")
  expect(actual_files).to_not include file_name
end


Then /^I (still|again) have an uncommitted file with name: "([^"]+)" and content: "([^"]+)"?$/ do |_, expected_name, expected_content|
  actual_files = run("git status --porcelain | awk '{print $2}'")[:out].split("\n")
  expect(actual_files).to eql [expected_name]

  # Verify the file content
  expect(IO.read expected_name).to eql expected_content
end

