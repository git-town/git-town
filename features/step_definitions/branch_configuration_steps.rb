# frozen_string_literal: true

Given(/^Git Town is aware of this branch hierarchy$/) do |table|
  table.hashes.each do |row|
    set_parent_branch branch: row['BRANCH'], parent: row['PARENT']
  end
end


Given(/^Git Town has no branch hierarchy information for "(.*?)"$/) do |branch_names|
  Kappamaki.from_sentence(branch_names).each do |branch_name|
    run_shell_command "git config --unset git-town-branch.#{branch_name}.parent"
  end
end


Then(/^(?:Git Town|it) is now aware of this branch hierarchy$/) do |table|
  table.diff! configured_branch_hierarchy_information.table
end


Then(/^Git Town has no branch hierarchy information$/) do
  expect(configured_branch_hierarchy_information(ignore_errors: true)).to be_empty
end


Then(/^I am not prompted for any parent branches$/) do
  expect(unformatted_last_run_output).not_to include 'Please specify the parent branch of'
end


Then(/^my branch hierarchy metadata is unchanged$/) do
  expect(@branch_hierarchy_metadata.table).to eql configured_branch_hierarchy_information.table
end
