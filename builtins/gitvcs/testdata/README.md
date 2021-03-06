# Integration Test Deploy Keys

The keys in this directory are configured as read-only [deploy
keys](https://docs.github.com/en/developers/overview/managing-deploy-keys#deploy-keys)
on the `gritcli/grit` repo. They are used for testing Grit's ability to clone
Git repositories.

It is safe to include the private key in the repo as it only grants read-only
access, and the `gritcli/grit` repo is already entirely public.

The passphrase for the `deploy-key-with-passphrase` key is `passphrase`.

The keys were generated with the following commands:

```console
ssh-keygen -t ed25519 -C "deploy-key-no-passphrase" -f deploy-key-no-passphrase -N ""
ssh-keygen -t ed25519 -C "deploy-key-with-passphrase" -f deploy-key-with-passphrase -N "passphrase"
```
