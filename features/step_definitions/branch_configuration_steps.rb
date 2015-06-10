Given(/^Git Town has no branch hierarchy information for "(.*?)"$/) do |branch_names|
  Kappamaki.from_sentence(branch_names).each do |branch_name|
    run_shell_command "git config --unset git-town.branches.parent.#{branch_name}"
    run_shell_command "git config --unset git-town.branches.parents.#{branch_name}"
  end
end

