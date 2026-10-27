# Provider SFTPGo

`provider-sftpgo` is a [Crossplane](https://crossplane.io/) provider for [SFTPGo](https://sftpgo.com/),
generated with [Upjet](https://github.com/crossplane/upjet) from
[drakkan/terraform-provider-sftpgo](https://github.com/drakkan/terraform-provider-sftpgo).

It currently exposes one managed resource, `User` (`user.sftpgo.m.zoispag.github.io`, plus a cluster-scoped
`user.sftpgo.zoispag.github.io`).

## How it works

The upstream Terraform provider is built on the plugin framework, so it is linked into the provider
binary and runs in-process: the image ships no Terraform CLI or provider plugin.

Only the S3 and local filesystem backends are exposed on `User`. The other backends (Azure Blob,
GCS, crypt, FTP, HTTP, SFTP) are pruned in `config/*/user/config.go`, together with the filesystem
SFTPGo reports back for virtual folders, so their credentials don't end up in the CRD. Upstream
doesn't mark the S3 `session_token` as sensitive; this provider does, so it is a secret reference
like `access_secret`. Open an issue if you need one of the other backends.

The user's password is read from a Secret through `passwordSecretRef`. The write-only
`password_wo` variant is not exposed.

## Configuration

The ProviderConfig secret takes `host` plus either `api_key`, or `username` and `password` of an
SFTPGo admin, the same settings as the Terraform provider block:

```json
{"host": "https://sftp.example.com", "api_key": "<sftpgo admin API key>"}
```

See `examples/namespaced/` for a `ClusterProviderConfig` and a `User`. Prefer the
`ClusterProviderConfig`: a namespaced `ProviderConfig` is looked up in the `User`'s namespace, and
its credentials secret must live in that same namespace.

The Terraform provider's `edition` and `headers` settings aren't passed through yet.

Leave `Delete` out of `managementPolicies` if removing the object should keep the SFTPGo account.

## Upgrading the Terraform provider

Bump `TERRAFORM_PROVIDER_VERSION` in the Makefile and the `github.com/drakkan/terraform-provider-sftpgo`
module in `go.mod` together, then run `make generate`.

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/zoispag/provider-sftpgo/issues).
