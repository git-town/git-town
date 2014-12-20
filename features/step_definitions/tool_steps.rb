Given(/^I have (.+?) installed$/) do |tools|
  names = Kappamaki.from_sentence(tools)
  names = [] if names == %w(nothing)
  update_installed_tools names
end
