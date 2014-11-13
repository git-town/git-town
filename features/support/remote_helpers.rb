def remote_url(name)
  output_of "git remote -v | grep '#{name}.*fetch' | awk '{print $2}'"
end
