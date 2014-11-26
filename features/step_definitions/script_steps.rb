Then(/^there (?:is an?|are) (.+) scripts? for "(.+?)"$/) do |actions, command|
  Kappamaki.from_sentence(actions).each do |action|
    expect(script_exists? command, action).to be_truthy, "#{script_path command, action} should exist"
  end
end


Then(/^there (?:is|are) no (.+) scripts? for "(.+?)" anymore$/) do |actions, command|
  Kappamaki.from_sentence(actions).each do |action|
    expect(script_exists? command, action).to be_falsy, "#{script_path command, action} should not exist"
  end
end
