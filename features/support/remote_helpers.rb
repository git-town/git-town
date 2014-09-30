def remote_url(name)
  run("git remote -v | grep '#{name}.*fetch' | awk '{print $2}'")[:out]
end
