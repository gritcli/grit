# The "daemon" block configures the Grit daemon and is also used by the CLI to
# determine how to communicate with the daemon.
#
# This block is optional, but if provided it may only be present in a single
# file within the configuration directory.
daemon {
  # The "socket" attribute is the path to the Unix socket used for
  # communication between the Grit CLI and the daemon. It defaults to
  # "~/grit/daemon.socket".
  socket = "/path/to/socket"
}

# The "clones" block configures Grit behaves with working with local clones of
# repositories by default.
#
# This block is optional, but if provided it may only be present in a single
# file within the configuration directory.
clones {
  # The "dir" attribute is the path to a directory into which local clones are
  # placed. It defaults to "~/grit".
  dir = "/path/to/clones"
}

# The "git" block configures how Grit's default behavior when working with Git
# repositories.
#
# This block is optional, but if provided it may only be present in a single
# file within the configuration directory.
git {
  # The "prefer_http" attribute instructs Grit to use the HTTP protocol for
  # Git operations whenever available. Otherwise, Grit prefers SSH. It
  # defaults to false.
  prefer_http = false

  # The "ssh_key" block explicitly defines an SSH key to use for Git
  # operations.
  #
  # This block is optional. By default Grit uses the system's SSH agent to
  # authenticate when using the SSH protocol.
  ssh_key {
    # The "file" attribute is the path to the SSH private key PEM file.
    file = "/path/to/key.pem"

    # The "passphrase" attribute is the passphrase to use to decrypt the
    # private key, if required. If it is empty it is assumed that the private
    # key is not encrypted.
    passphrase = "<passphrase>"
  }
}

# A "source" block defines a repository source that hosts the repositories that
# may be cloned by Grit.
#
# A configuration may contain any number of "source" blocks.
#
# The first parameter is the source "name", which must be unique to this source.
# Names are case-insensitive and may contain alpha-numeric characters and
# underscores.
#
# The second parameter is the "driver", which determines how Grit communicates
# with the remote server. Configuration examples for each of the built-in
# drivers are given later in the file.
source "example_source" "some_driver" {
  # The "enabled" parameter controls whether the source is enabled for
  # operations such as cloning new repositories. It defaults to true.
  enabled = true

  # The "clones" block configures how Grit behaves when working with local
  # clones of repositories that were obtained from this soure.
  #
  # This block is optional, by default Grit basis its behavior of the top-level
  # "git" block.
  clones {
    # The "dir" attribute is the path to a directory into which local clones are
    # placed.
    #
    # By default, clones from this source are kept in a sub-directory of the
    # directory configured by the top-level "git" block. The sub-directory is
    # named the same as the source.
    #
    # For example, if the top-level "git" block uses a "dir" of "~/clones", the
    # default for this source would be "~/clones/example_source".
    dir = "/path/to/somewhere/else"
  }

  # Each source may contain additional options that are specific to the chosen
  # driver.
  some_driver_specific_option = "<value>"
}

# This source demonstrates the configuration options that are unique to the
# built-in "github" driver which can be used for sources that use GitHub.com or
# GitHub Enterprise Server.
source "example_github_source" "github" {
  # The "domain" attribute is the domain name where the GitHub server is
  # located. It defaults to "github.com".
  domain = "code.example.org"

  # The "token" attribute is the GitHub PAT (personal access token) used to
  # authenticate against the GitHub API.
  #
  # By default the "github" driver works without authenticating, though with a
  # reduced feature set.
  token = "<github auth token>"

  # The "git" block configures how Grit behaves when working with Git
  # repositories that were obtained from this source.
  #
  # It may contain the same attributes as the top-level "git" block. Any value
  # specified here overrides the value specified in the top-level "git" block.
  git {
    # ...
  }
}

# TODO: example sources for other drivers
