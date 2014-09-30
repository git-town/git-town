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
