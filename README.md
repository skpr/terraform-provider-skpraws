# Terraform provider skpraws

This provider implements aws endpoints that are currently not available in upstream repositories.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the command:

```shell
make install
```

## Using the provider

See the `examples/` for using the provider.

## Developing the Provider

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

Update `~/.terraformrc` to use the provider locally.

```
provider_installation {
  dev_overrides {
    "skpr/skpraws" = "/my/path/to/go/bin"
  }
  # Install other providers normally from the registry
  direct {}
}
```

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.
