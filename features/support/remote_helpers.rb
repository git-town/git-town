# Returns the base URL for the given domain
def base_url domain
  case domain
  when 'Bitbucket' then 'https://bitbucket.org'
  when 'GitHub' then 'https://github.com'
  else fail "Unknown domain: #{domain}"
  end
end


# Returns the remote git URL for domain / protocol
def git_url domain, protocol
  "#{git_url_prefix domain, protocol}Originate/git-town.git"
end


# Returns the remote git URL prefix for the given domain and protocol
def git_url_prefix domain, protocol
  case [domain, protocol]
  when %w(Bitbucket HTTPS) then 'https://username@bitbucket.org/'
  when %w(Bitbucket SSH) then 'git@bitbucket.org:'
  when %w(GitHub HTTPS) then 'https://github.com/'
  when %w(GitHub SSH) then 'git@github.com:'
  else fail "Unknown domain/protocol pairing: #{domain}/#{protocol}"
  end
end


# Returns the remote URL for a new pull request for the given domain and branch
def pull_request_url domain, branch_name
  if domain == 'Bitbucket'
    sha = recent_commit_shas(1).join('')[0, 12]
    "https://bitbucket.org/Originate/git-town/pull-request/new?source=Originate%2Fgit-town%3A#{sha}%3A#{branch_name}"
  elsif domain == 'GitHub'
    "https://github.com/Originate/git-town/compare/#{branch_name}?expand=1"
  else
    fail "Unknown domain: #{domain}"
  end
end


# Returns the remote URL for the homepage of the given domain
def repository_homepage_url domain
  "#{base_url domain}/Originate/git-town"
end
