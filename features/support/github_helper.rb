require 'net/http'
require 'uri'

def github_rails_fork
  $github_rails_fork ||= (
    uri = URI('https://api.github.com/repos/rails/rails/forks?page=1&per_page=1')
    response = Net::HTTP.get_response(uri)
    forks = JSON.parse(response.body)
    forks[0]
  )
end

def github_check_rate_limit!
  uri = URI('https://api.github.com/rate_limit')
  response = Net::HTTP.get_response(uri)
  rate_limit = JSON.parse(response.body)
  remaining = rate_limit['rate']['remaining']
  reset = Time.at(rate_limit['rate']['reset']).strftime("%I:%M:%S %P")

  # Grabbing the rails forks + 4 tests hit the API
  if remaining < 5
    raise "Github API rate limit reached - will reset at #{reset}"
  end
end
