Given(/^I have an uncommitted file with name: "(.*?)" and content: "(.*?)"$/) do |name, content|
  IO.write name, content
end


Given(/^I resolve the conflict in "(.+?)"$/) do |file_name|
  IO.write file_name, 'resolved content'
  run "git add #{file_name}"
end




Then(/^(?:now I|I still) have the following committed files$/) do |files_data|
  files_data.map_headers!(&:downcase)

  # Get all expected files
  expected_files = files_data.hashes.map do |expected_file|
    symbolize_keys_deep! expected_file
    filenames = expected_file.delete :files
    Kappamaki.from_sentence(filenames).map do |filename|
      result = expected_file.clone
      result[:name] = filename
      result
    end
  end.flatten

  # Get all existing files in all branches
  actual_files = []
  without_open_changes do
    existing_local_branches.each do |branch|
      on_branch branch do
        existing_files.each do |file|
          if file != 'uncommitted'
            actual_files << { branch: branch, name: file, content: IO.read(file) }
          end
        end
      end
    end
  end

  # Remove the keys that are not used in the expected data
  used_keys = expected_files[0].keys
  actual_files.each do |actual_file|
    actual_file.keys.each do |key|
      actual_file.delete key unless used_keys.include? key
    end
  end

  expect(actual_files).to match_array expected_files
end


Then(/^I don't have an uncommitted file with name: "(.*?)"$/) do |file_name|
  expect(uncommitted_files).to_not include file_name
end


Then(/^I (?:still|again) have an uncommitted file with name: "([^"]+)" and content: "([^"]+)"$/) do |file_name, content|
  expect(uncommitted_files).to eql [file_name]
  expect(IO.read file_name).to eql content
end
