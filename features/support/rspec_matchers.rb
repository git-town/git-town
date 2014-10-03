RSpec::Matchers.define :be_success do
  match do |actual|
    actual[:status] == 0
  end
  failure_message_for_should do |actual|
    "Operation was not successful! Out: '#{actual.out}', Err: '#{actual.err}'"
  end
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
