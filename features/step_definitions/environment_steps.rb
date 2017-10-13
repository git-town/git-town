# frozen_string_literal: true

Given(/^my workspace is currently not in a git repository$/) do
  FileUtils.rm_rf '.git'
end
