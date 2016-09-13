Given(/^I'm currently not in a git repository$/) do
  FileUtils.rm_rf '.git'
end
