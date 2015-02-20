# rubocop:disable MethodLength
def verify_error message
  @error_expected = true
  expect(@last_run_result.error).to be_truthy

  actual = unformatted_last_run_output
  expected = message.gsub(/\W/, '')
  expect(actual).to include(expected), %(
    ACTUAL
    ***************************************************
    #{actual.dump}
    ***************************************************
    EXPECTED TO INCLUDE
    ***************************************************
    #{expected.dump}
    ***************************************************
  ).gsub(/^ {4}/, '')
end
# rubocop:enable MethodLength
