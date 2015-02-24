# Pauses the spec to let the developer inspect the test repos.
#
# This step is for debugging purposes only.
# Please keep it around.
Then(/^show me my repo$/) do
  p "Developer repo: #{repository_path :developer}"
  p "Coworker repo: #{repository_path :coworker}"
  p 'Press ENTER to continue'
  STDIN.gets
end
