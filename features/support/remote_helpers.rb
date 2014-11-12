def remote_url(name)
  get_output("git remote -v | grep '#{name}.*fetch' | awk '{print $2}'")
end
