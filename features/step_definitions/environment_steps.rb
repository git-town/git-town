Given(/^I'm currently not in a Git  repository$/) do
  FileUtils.rm_rf '.git'
end
