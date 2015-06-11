require 'cgi'


# Returns the base URL for the given domain
def base_url domain
  case domain
  when 'Bitbucket' then 'https://bitbucket.org'
  when 'GitHub' then 'https://github.com'
  else fail "Unknown domain: #{domain}"
  end
end


# Returns the remote URL for a new pull request for the given domain and branch
def pull_request_url domain, branch, parent_branch, repo
  case domain
  when 'Bitbucket'
    sha = recent_commit_shas(1).join('')[0, 12]
    "https://bitbucket.org/#{repo}/pull-request/new?source=#{CGI.escape repo}%3A#{sha}%3A#{branch}"
  when 'GitHub'
    "https://github.com/#{repo}/compare/#{parent_branch}...#{branch}?expand=1"
  else
    fail "Unknown domain: #{domain}"
  end
end


# Returns the remote URL for the homepage of the given domain
def repository_homepage_url domain, repository
  "#{base_url domain}/#{repository}"
end
