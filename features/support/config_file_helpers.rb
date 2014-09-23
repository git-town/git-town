def config_file_path version
  case version
    when 'old' then '.main_branch_name'
    when 'new' then '.gittownrc'
    else raise "Unknown config file version: '#{version}'"
  end
end

