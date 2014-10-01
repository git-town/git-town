# Helper methods for working with GitHub.


# Parses the remote url extracting the GitHub repository
# Returns empty string if not a GitHub repository
function github_parse_repository {
  echo `echo "$1" | sed -n "s/.*github.com[/:]\(.*\).git/\1/p"`
}

# Queries GitHub for the upstream url for the given repository and protocol
function github_upstream_url {
  local repository=$1
  local protocol=$2
  local parse_json="""
import json, sys
obj = json.load(sys.stdin)
if 'parent' in obj:
  print obj['parent']['$protocol']
"""

  ensure_tool_installed python
  echo `curl -s https://api.github.com/repos/$repository | python -c "$parse_json"`
}
