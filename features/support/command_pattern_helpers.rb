def mock_script_create
  file_contents = "#!/usr/bin/env bash\n"
  file_contents << "source \"#{File.dirname(__FILE__)}/../../src/helpers/helpers.sh\" --bypass-automatic-configuration --bypass-environment-checks\n"
  file_contents <<
    <<-END.gsub(/^ {6}/, '')
      function preconditions {
        # PRECONDITIONS GO HERE
      }

      function steps {
        # STEPS GO HERE
      }

      run "$@"
    END

  FileUtils.mkdir_p MOCK_SCRIPT_DIRECTORY
  File.open(MOCK_SCRIPT_FILENAME, 'w') do |file|
    file.write file_contents
  end
  FileUtils.chmod 'u+x', MOCK_SCRIPT_FILENAME
end


def mock_script_update function, commands
  mock_script_contents = File.read(MOCK_SCRIPT_FILENAME)
  mock_script_contents.gsub!(/# #{function.upcase} GO HERE/, commands.join("\n"))

  File.open(MOCK_SCRIPT_FILENAME, 'w') {|file| file.puts mock_script_contents }
end