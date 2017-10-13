# frozen_string_literal: true

Given(/^my workspace is currently not a Git repository$/) do
  FileUtils.rm_rf '.git'
end
