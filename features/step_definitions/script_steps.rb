Then(/^there (?:is an?|are) (.+) scripts? for "(.+)"$/) do |actions, command|
  Kappamaki.from_sentence(actions).map do |action|
    expect(script_exists? command, action).to be_truthy, "#{action} script for #{command} should exist"
  end
end


Then(/^there (?:is|are) no (.+) scripts? for "(.+)" anymore$/) do |actions, command|
  Kappamaki.from_sentence(actions).map do |action|
    expect(script_exists? command, action).to be_falsy, "#{action} script for #{command} should not exist"
  end
end
