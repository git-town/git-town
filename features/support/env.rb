require 'kappamaki'
require 'open4'
require 'rspec'
require 'active_support/all'

# The files to ignore when checking files
IGNORED_FILES = %w[ tags ]

Before do
  Dir.chdir repositiory_base

  # Create remote repository
  create_repository remote_repository_path

  # Create the local repository
  clone_repository remote_repository_path, local_repository_path
  at_path local_repository_path do

    # Create the master branch
    run 'touch .gitignore ; git add .gitignore ; git commit -m "Initial commit"; git push -u origin master'

    # Create the main branch
    run 'git checkout -b main master ; git push -u origin main'
    run 'git config git-town.main-branch-name main'
  end

  # Create the coworker repository
  clone_repository remote_repository_path, coworker_repository_path
  at_path coworker_repository_path do
    run 'git checkout main'
    run 'git config git-town.main-branch-name main'
  end

  Dir.chdir local_repository_path
end


Before('@github_query') do
  github_check_rate_limit!
end


at_exit do
  Dir.chdir repositiory_base
  delete_repository remote_repository_path
  delete_repository local_repository_path
  delete_repository coworker_repository_path
end


RSpec::Matchers.define :match_commits do |expected|
  match do |actual|
    (expected - actual).empty?
  end

  failure_message_for_should do |actual|
    result = ""
    expected.sort! {|x, y| commit_sorted(x) <=> commit_sorted(y) }
    actual.sort! {|x, y| commit_sorted(x) <=> commit_sorted(y) }

    result << "\nEXPECTED VALUES\n"
    expected.each do |commit|
      result << commit_to_s(commit)
    end

    result << "\nACTUAL VALUES\n"
    actual.each do |commit|
      result << commit_to_s(commit)
    end

    result << "\nCOMMON COMMITS\n"
    common_commits = expected & actual
    common_commits.each do |commit|
      result << commit_to_s(commit)
    end

    expected_but_not_present = expected - actual
    unless expected_but_not_present.empty?
      result << "\nEXPECTED BUT NOT PRESENT COMMITS:\n"
      expected_but_not_present.each do |commit|
      result << commit_to_s(commit)
      end
    end

    present_but_not_expected = actual - expected
    unless present_but_not_expected.empty?
      result << "\nPRESENT BUT NOT EXPECTED COMMITS:\n"
      present_but_not_expected.each do |commit|
      result << commit_to_s(commit)
      end
    end

    result + "\n"
  end

  def commit_to_s commit
    "#{commit[:branch]} branch (#{commit[:location]}): '#{commit[:message]}' with #{commit[:files]}\n"
  end

  def commit_sorted commit
    "#{commit[:message]} #{commit[:branch]} #{commit[:location]} #{commit[:files]}"
  end
end


