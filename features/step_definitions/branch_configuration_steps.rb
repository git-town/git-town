Given(/^Git Town is aware of this branch hierarchy$/) do |table|
  table.hashes.each do |row|
    set_parent_branch branch: row['BRANCH'], parent: row['PARENT']
  end
end


Given(/^Git Town has no branch hierarchy information for "(.*?)"$/) do |branch_names|
  Kappamaki.from_sentence(branch_names).each do |branch_name|
    run_shell_command "git config --unset git-town.branches.parent.#{branch_name}"
    run_shell_command "git config --unset git-town.branches.ancestors.#{branch_name}"
  end
end


Then(/^Git Town is now aware of this branch hierarchy$/) do |table|
  table.diff! configured_branch_hierarchy_information.table
end


Then(/^Git Town is not aware of any branch hierarchy$/) do
  expect(configured_branch_hierarchy_information(ignore_errors: true).table.size).to eq 1
end


Then(/^my branch hierarchy metadata is unchanged$/) do
  expect(@branch_hierarchy_metadata.table).to eql configured_branch_hierarchy_information.table
end
