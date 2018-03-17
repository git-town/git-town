# Website

The Git Town website can be found at <http://www.git-town.com>.
It is hosted with [GitHub pages](https://pages.github.com/)

## Requirements

* [Node.js](https://nodejs.org)
  * on OS X best installed via [Homebrew](http://brew.sh)
* [Harp](http://harpjs.com)
  * `npm install -g harp`

## Local Development

* `cd website`
* `harp server`
* go to [localhost:9000](http://localhost:9000)

## Deployment

* make sure your changes to be deployed are all in `master`
* `make deploy`
