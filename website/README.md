# Website Development

The Git Town website is at [git-town.com](https://www.git-town.com).

### local development

The source code for the website is in the [website](../website/) folder. This
folder contains its own [Makefile](../website/Makefile) for activities related
to working on the website. Run `make setup` to download the necessary tooling,
`make serve` to start a local development server, and `make docs` to test the
website.

### production environment

The website runs on [Netlify](https://www.netlify.com). It auto-updates on
changes to the `main` branch. The Netlify configuration is in
[netlify.toml](../netlify.toml)
