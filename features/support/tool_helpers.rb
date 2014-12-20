def update_installed_tools tools
  File.open(TOOLS_INSTALLED_FILENAME, 'w') do |file|
    file.write tools.join("\n") + "\n"
  end
end
