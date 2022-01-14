# Website Development Guide

The Git Town website is at [git-town.com](https://www.git-town.com). The website
runs on [Netlify](https://www.netlify.com). It auto-updates on changes to the
`main` branch.

To run the website locally, install
[mdBook](https://github.com/rust-lang/mdBook). To test the website, install
[Node.js](https://nodejs.org) and [Yarn](https://yarnpkg.com), then run
<code textrun="verify-make-command">make setup</code>.

The source code is in the [website](../../website/) folder.

- run a local dev server: `make website-dev`
- test that the website compiles: `make website-build`
