# Website Development Guide

The Git Town website is at https://www.git-town.com.

## setup

Install [mdBook](https://github.com/rust-lang/mdBook) to transpile the website
CMS. Install

- install [Node.js](https://nodejs.org) and [Yarn](https://yarnpkg.com)

## local development

- run a local dev server: `make website-dev`
- test that the website compiles: `make website-build`

## deployment

The website runs on [Netlify](https://www.netlify.com). It auto-updates on
changes to the `main` branch.
