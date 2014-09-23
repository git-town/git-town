RSpec::Matchers.define :be_success do
  match do |actual|
    actual[:status] == 0
  end
  failure_message_for_should do |actual|
    "Operation was not successful! Out: '#{actual.out}', Err: '#{actual.err}'"
  end
end

