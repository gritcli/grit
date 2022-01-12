# Grit

[![Build Status](https://github.com/gritcli/grit/workflows/CI/badge.svg)](https://github.com/gritcli/grit/actions?workflow=CI)
[![Code Coverage](https://img.shields.io/codecov/c/github/gritcli/grit/main.svg)](https://codecov.io/github/gritcli/grit)
[![Latest Version](https://img.shields.io/github/tag/gritcli/grit.svg?label=semver)](https://semver.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/gritcli/grit)](https://goreportcard.com/report/github.com/gritcli/grit)

Grit is a tool for managing local clones of Git (and other) repositories from
your terminal.

It allows quickly cloning repositories and switching the current directory to
local clones based on a (partial) repository name, with shell auto-completion.

Grit integrates with cloud services such as GitHub and can also be configured to
work with self-hosted version control systems.

> **NOTE:** This repository contains Grit version 2 which is under active
> development and is by no means ready for use. Grit version 1 is available at
> [`jmalloc/grit`](https://github.com/jmalloc/grit).
>
> **Everything below this point in the README is likely not yet functional or
> subject to change.**

---

```
TODO: asciinema recording of "clone", "go" and "cd" commands
```

## Installation

```
TODO: use [n]fpm to publish more "acceptable" linux packages
```

### Install via Homebrew

Grit can be installed via [Homebrew](https://brew.sh) on **macOS** or **Linux**
by running the commands below:

```bash
brew tap gritcli/grit
brew install grit
```

Alternatively, install using [Homebrew
Bundle](https://github.com/Homebrew/homebrew-bundle) by adding the following
lines to your `Brewfile`:

```Brewfile
tap "gritcli/grit"
brew "grit"
```

### Install manually

1. Download the Grit archive for your platform from the [latest release](https://github.com/gritcli/grit/releases/latest).
2. Launch the `gritd` in the background when your system starts
3. Put the `grit` (cli) executable in your path

The `gritd` daemon is designed to be run as your regular system user; it does
not require elevated privileges.

## Getting started

```
TODO
```

## Concepts

For historical reasons (and to avoid overly abstract language) Grit largely uses
Git's terminology for version control concepts. For example, a local working
copy of a repository is referred to as a "clone".

### Source

A source is some remote system that hosts VCS repositories that can be cloned by
Grit.

Grit can be configured to consume any number of sources. Each source uses a
specific [driver](#driver) to communicate with the remote system.

Each source is identified by a unique name. There are several built-in sources:

- [ ] `bitbucket` for repositories hosted on [BitBucket Cloud](https://bitbucket.org/product/)
- [x] `github` for repositories hosted on [GitHub.com](https://github.com)
- [ ] `gitlab` for repositories hosted on [GitLab.com](https://gitlab.com/explore)

Additionally, user-defined sources can be configured to consume repositories
from self-hosted VCS systems.

### Driver

A driver integrates Grit with a specific kind of [source](#source). It
encapsulates all of the communication with the remote source, such as API calls
and VCS operations.

Grit ships with several built-in drivers:

- [ ] `bitbucket` for [BitBucket Cloud, BitBucket Server and BitBucket Data Center](https://bitbucket.org/product/guides/getting-started/overview#bitbucket-software-hosting-options)
- [ ] `gitea` for [Gitea](https://gitea.io)
- [x] `github` for [GitHub.com](https://github.com) and [GitHub Enterprise Server](https://docs.github.com/en/get-started/signing-up-for-github/setting-up-a-trial-of-github-enterprise-server)
- [ ] `gitlab` for [GitLab.com](https://gitlab.com/explore) and [Self-managed GitLab](https://about.gitlab.com/install/)
- [ ] `gogs` for [Gogs](https://gogs.io)

Additionally, custom drivers can be implemented via plugins. There is no
requirement that a driver use Git as its underlying VCS.

## Configuration

Grit works out-of-the-box with zero configuration, however more powerful
features can be enabled with some configuration.

Grit configuration files are written in
[HCL](https://github.com/hashicorp/hcl#why). Grit loads all `.hcl` files in the
`~/.config/grit` directory by default. Files that begin with an underscore or
dot are ignored.

```
TODO: provide a guide for things that the user will most likely want to configure:

- Authentication
- Custom sources
```

The [`config-reference.hcl`](config-reference.hcl) file demonstrates all of the
available configuration options and their default values.

## Migrating from Grit version 1

```
TODO
```

## History and rationale

I spend most of my day working with Git. Many of the repositories are hosted on
GitHub.com, but many more are in my employer's private GitHub Enterprise and
BitBucket installations.

Keeping track of hundreds of clones can be a little tedious, so some time back
I adopted a basic directory naming convention and wrote some shell scripts to
handle cloning in a consistent way.

This worked well for a while, until the list of places I needed to clone from
increased further, and I started working more heavily in [Go](http://golang.org),
which, at the time, placed it's [own requirements](https://github.com/golang/go/wiki/GOPATH)
on the location of your Git clones.

Grit is the logical evolution of those original scripts into a standalone
project that clones from multiple Git sources.

Grit v1 was hacked together, there were no tests, and there are other more
general solutions for navigating your filesystem; but it worked for me and my
colleagues.

Grit v2 is an attempt to address some of the feature requests from my colleagues
and to make Grit a more pleasant project to maintain, with better internal
abstractions, good test coverage and an eye towards extensibility.
