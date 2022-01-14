# Website Development Guide

The Git Town website is at [git-town.com](https://www.git-town.com). The website
runs on [Netlify](https://www.netlify.com). It auto-updates on changes to the
`main` branch.

The source code is in the [website](../../website/) folder. This folder contains
its own [Makefile](../../website/Makefile). To run the website locally, install
[mdBook](https://github.com/rust-lang/mdBook) and run `make serve`. To test the
website, install [Node.js](https://nodejs.org) and [Yarn](https://yarnpkg.com),
then run
<code textrun="verify-make-command">make setup</code> to install dependencies,
then <code textrun="verify-make-command">make docs</code> to verify the link
structure.
