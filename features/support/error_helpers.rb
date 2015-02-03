def verify_error message
  @error_expected = true

  expect(@last_run_result.error).to be_truthy
  expect(unformatted_last_run_output.strip).to include(message), %(
    ACTUAL
    ***************************************************
    #{@last_run_result.out.gsub '\n', "\n"}
    ***************************************************
    EXPECTED TO INCLUDE
    ***************************************************
    #{message.gsub '\n', "\n"}
    ***************************************************
  ).gsub(/^ {4}/, '')
end
