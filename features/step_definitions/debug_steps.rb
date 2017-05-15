# frozen_string_literal: true
# These steps are for debugging purposes only.
# Please keep them around.


# Displays the output of the given command.
Then(/^show me the output of `(.+?)`$/) do |command|
  puts output_of command
end


# Pauses the spec to let the developer inspect the test repos.
Then(/^show me my repo$/) do
  p "Developer repo: #{repository_path :developer}"
  p "Coworker repo: #{repository_path :coworker}"
  p 'Press ENTER to continue'
  STDIN.gets
end
