# Integration Test Deploy Keys

The keys in this directory are configured as read-only [deploy
keys](https://docs.github.com/en/developers/overview/managing-deploy-keys#deploy-keys)
on the `gritcli/grit` repo. They are used for testing Grit's integration with GitHub, specifically for
cloning repositories.

It is safe to include the private key in the repo as it only grants read-only
access, and the `gritcli/grit` repo is already entirely public.

The passphrase for the `deploy-key-with-passphrase` key is `passphrase`.
