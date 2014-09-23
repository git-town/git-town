Given /^I don't have a configuration file$/ do
  File.delete '.gittownrc'
end


Given /^I have an? (old|new) configuration file with (.*?)$/ do |version, data|
  File.delete '.gittownrc'
  data = Kappamaki.attributes_from_sentence data
  IO.write config_file_path(version), data[:'main branch']
end


Then /^I end up with a (old|new) configuration file with (.*?)$/ do |version, data|
  file_name = config_file_path version
  data = Kappamaki.attributes_from_sentence data
  expect(File.exist? file_name).to be_truthy, "Could not find file #{file_name}"
  expect(IO.read file_name).to eq data[:'main branch']
end


Then /^I don't have an old configuration file anymore$/ do
  expect(File.exist? config_file_path('old')).to be_falsy
end
