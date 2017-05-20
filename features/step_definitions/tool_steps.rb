# frozen_string_literal: true
Given(/^I have "(.+?)" installed$/) do |tool|
  @tool = tool
  update_installed_tools [tool]
end


Given(/^I have Git "(\d+)\.(\d+)\.(\d+)" installed$/) do |major, minor, patch|
  gitPath = File.join(@temporary_shell_overrides_directory, 'git')
  IO.write gitPath, <<~HEREDOC
    #!/usr/bin/env bash
    echo "git version #{major}.#{minor}.#{patch}"
  HEREDOC
  FileUtils.chmod "u+x", gitPath, :verbose => true
end


Given(/^I have no command that opens browsers installed$/) do
  update_installed_tools []
end
