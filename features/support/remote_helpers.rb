# Returns the remote URL for domain / protocol
#   domain is Github or Bitbucket
#   protocol is HTTPS or SSH
def remote_url domain, protocol
  if protocol == 'HTTPS'
    host = domain == 'Github' ? 'github.com' : 'username@bitbucket.org'
    "https://#{host}/Originate/git-town.git"
  else
    host = domain == 'Github' ? 'github.com' : 'bitbucket.org'
    "git@#{host}:Originate/git-town.git"
  end
end

# Returns the remote URL for a new pull request for the given domain and branch
#   domain is Github or Bitbucket
def remote_pull_request_url domain, branch_name
  if domain == 'Github'
    "https://github.com/Originate/git-town/compare/#{branch_name}?expand=1"
  else
    sha = recent_commit_shas(1).join('')[0, 12]
    "https://bitbucket.org/Originate/git-town/pull-request/new?source=Originate%2Fgit-town%3A#{sha}%3A#{branch_name}"
  end
end
