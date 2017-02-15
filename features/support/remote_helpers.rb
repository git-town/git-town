require 'cgi'


# Returns the base URL for the given domain
def base_url domain
  case domain
  when 'Bitbucket' then 'https://bitbucket.org'
  when 'GitHub' then 'https://github.com'
  when 'GitLab' then 'https://gitlab.com'
  else fail "Unknown domain: #{domain}"
  end
end


# Returns the URL for making pull requests on Bitbucket
def bitbucket_pull_request_url branch:, parent_branch:, repo:
  sha = recent_commit_shas(1).join('')[0, 12] # TODO: update to have the branch as an argument
  source = CGI.escape "#{repo}:#{sha}:#{branch}"
  dest = CGI.escape "#{repo}::#{parent_branch}"
  "https://bitbucket.org/#{repo}/pull-request/new?source=#{source}&dest=#{dest}"
end


# Returns the URL for making pull requests on GitHub
def github_pull_request_url branch:, parent_branch: nil, repo:
  to_compare = parent_branch ? "#{parent_branch}...#{branch}" : branch
  "https://github.com/#{repo}/compare/#{to_compare}?expand=1"
end

# Returns the URL for making pull requests on GitLab
def gitlab_pull_request_url branch:, parent_branch: nil, repo:
  to_compare = parent_branch ? "#{parent_branch}...#{branch}" : branch
  "https://gitlab.com/#{repo}/compare/#{to_compare}?expand=1"
end


# Returns the remote URL for a new pull request for the given domain and branch
def pull_request_url domain:, branch:, parent_branch: nil, repo:
  send "#{domain.downcase}_pull_request_url",
       branch: branch, parent_branch: parent_branch, repo: repo
end


# Returns the remote URL for the homepage of the given domain
def repository_homepage_url domain, repository
  "#{base_url domain}/#{repository}"
end
