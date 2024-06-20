# Welcome 💖

Before anything else, thank you for taking some of your precious time to help this project move forward. ❤️

If you're new to open source and feeling a bit nervous 😳, we understand! We recommend watching [this excellent guide](https://egghead.io/talks/git-how-to-make-your-first-open-source-contribution)
to give you a grounding in some of the basic concepts. You could also watch [this talk](https://www.youtube.com/watch?v=tuSk6dMoTIs) from our very own wonderful [Marit van Dijk](https://github.com/mlvandijk) on her experiences contributing to Cucumber.

We want you to feel safe to make mistakes, and ask questions. If anything in this guide or anywhere else in the codebase doesn't make sense to you, please let us know! It's through your feedback that we can make this codebase more welcoming, so we'll be glad to hear thoughts.

You can chat with us in the [#committers-go](https://cucumberbdd.slack.com/archives/CA5NJPDJ4) channel in our [community Slack], or feel free to [raise an issue] if you're experiencing any friction trying make your contribution.

## Setup

To get your development environment set up, you'll need to [install Go]. We're currently using version 1.17 for development.

Once that's done, try running the tests:

    make test

If everything passes, you're ready to hack!

[install go]: https://golang.org/doc/install
[community Slack]: https://cucumber.io/community#slack
[raise an issue]: https://github.com/cucumber/godog/issues/new/choose

## Changing dependencies

If dependencies have changed, you will also need to update the _examples module. `go mod tidy` should be sufficient.