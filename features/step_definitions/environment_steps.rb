Given(/^I'm currently not in a git town-repository$/) do
  FileUtils.rm_rf '.git'
end
