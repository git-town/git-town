def verify_error message
  @error_expected = true
  expect(@last_run_result.error).to be_truthy

  expect(unformatted_last_run_output).to include(message), %(
    ACTUAL
    ***************************************************
    #{unformatted_last_run_output.dump}
    ***************************************************
    EXPECTED TO INCLUDE
    ***************************************************
    #{message.dump}
    ***************************************************
  ).gsub(/^ {4}/, '')
end
