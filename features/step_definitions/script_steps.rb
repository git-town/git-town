Then(/^there (?:is|are) no (.+) scripts? for "(.+)"$/) do |actions, command|
  Kappamaki.from_sentence(actions).map do |action|
    expect(script_exists? command, action).to be_falsy, "#{action} script for #{command} should not exist"
  end
end
