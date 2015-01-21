def verify_error method, message
  @error_expected = true

  expect(@last_run_result.error).to be_truthy
  expect(unformatted_last_run_output).to public_send(method, message), %(
    ACTUAL
    ***************************************************
    #{@last_run_result.out.gsub '\n', "\n"}
    ***************************************************
    EXPECTED TO #{method}
    ***************************************************
    #{message.gsub '\n', "\n"}
    ***************************************************
  ).gsub(/^ {4}/, '')
end
