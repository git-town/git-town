# Returns the base URL for the given domain
def base_url domain
  case domain
  when 'Bitbucket' then 'https://bitbucket.org'
  when 'GitHub' then 'https://github.com'
  else fail "Unknown domain: #{domain}"
  end
end


# Returns the remote git URL for domain / protocol
def git_url domain, protocol, suffix
  "#{git_url_prefix domain, protocol}Originate/git-town#{suffix}"
end


# Returns the remote git URL prefix for the given domain and protocol
def git_url_prefix domain, protocol
  case domain
  when 'Bitbucket' then git_url_prefix_bitbucket protocol
  when 'GitHub' then git_url_prefix_github protocol
  else fail "Unknown domain: #{domain}"
  end
end


# Returns the remote git URL prefix on Bitbucket for the given protocol
def git_url_prefix_bitbucket protocol
  case protocol
  when 'HTTP', 'HTTPS' then "#{protocol.downcase}://username@bitbucket.org/"
  when 'SSH' then 'git@bitbucket.org:'
  else fail "Unknown protocol: #{protocol}"
  end
end


# Returns the remote git URL prefix on GitHub for the given protocol
def git_url_prefix_github protocol
  case protocol
  when 'HTTP', 'HTTPS' then "#{protocol.downcase}://github.com/"
  when 'SSH' then 'git@github.com:'
  else fail "Unknown protocol: #{protocol}"
  end
end


# Returns the remote URL for a new pull request for the given domain and branch
def pull_request_url domain, branch_name
  case domain
  when 'Bitbucket'
    sha = recent_commit_shas(1).join('')[0, 12]
    "https://bitbucket.org/Originate/git-town/pull-request/new?source=Originate%2Fgit-town%3A#{sha}%3A#{branch_name}"
  when 'GitHub'
    "https://github.com/Originate/git-town/compare/#{branch_name}?expand=1"
  else
    fail "Unknown domain: #{domain}"
  end
end


# Returns the remote URL for the homepage of the given domain
def repository_homepage_url domain
  "#{base_url domain}/Originate/git-town"
end
