# Bob, the builder

[![Maintainability](https://api.codeclimate.com/v1/badges/09a3057a188cff0b2f1c/maintainability)](https://codeclimate.com/github/vitorfhc/bob/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/09a3057a188cff0b2f1c/test_coverage)](https://codeclimate.com/github/vitorfhc/bob/test_coverage)

> Bob is an automated multiple Docker images builder. It works by defining a simple YAML file and running its commands.

### Under development

This project is under heavy development.

Tasks:

- [x] Feature: build images
- [x] Feature: push images
- [x] Tests: pkg/helpers
- [x] Improvement: use [spf13/viper](https://github.com/spf13/viper) for better configuration
- [x] Feature: dependencies between images
- [ ] Tests: pkg/docker
- [ ] Feature: parallel building
- [ ] Feature: parallel pushing

## Installing

### Using Go

Run `go install` command as showed below:

```bash
$ export BOB_VERSION=latest # specify the version you want (example: v0.0.6)
$ go install github.com/vitorfhc/bob@$BOB_VERSION
```

### From source

To install the project from its source code you need to clone this repository and have `golang` installed.

```bash
$ git clone https://github.com/vitorfhc/bob
$ cd bob
$ go build -o bob
$ sudo mv bob /usr/local/bin
```

## Using bob

Running `bob` command gives you the following usage message:

```
Using this tool you may build and push several images in a monorepo.
All you need is a YAML file which has everything you need configured.

Examples:
  bob build

Usage:
  bob [command]

Available Commands:
  build       Builds Docker images
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  push        Pushed Docker images

Flags:
  -d, --debug         wether to print debug messages
  -h, --help          help for bob

Use "bob [command] --help" for more information about a command.
```

### Configuration File (bob.yaml)

### root schema

| field | type | description | default | required |
|----|---|---|---|---|
| .debug | boolean | if debug mode is on | false | |
| .username | string | username to be used on push | "" | |
| .password | string | password to be used on push | "" | |
| .images | list | list of images ([image schema](#image-schema)) | [] | |

### image schema

| field | type | description | default | required |
|----|---|---|---|---|
| .id | string | unique identifier for the image | random string | |
| .name | string | name of the image | | X |
| .tags | list of strings | image tags | [] (this defaults to a single tag `latest`) | |
| .context | string | build context | current dir | |
| .dockerfile | string | path to Dockerfile relative to context | Dockerfile | |
| .target | string | Dockerfile target | | |
| .buildArgs | map of strings | build args | | |
| .registry | string | image registry to push to | "" | |
| .needs | list of strings | strings referring to other images IDs that must be built before this one | [] | |

### Example

In the [examples folder](examples) you can find some ways to use `bob`.

Take a look at [this example's](examples/example_03/bob.yaml) YAML:

```yaml
images:
- name: alpine-dynamic
  tags:
  - "3.14-dynamic"
  context: .
  dockerfile: alpine/Dockerfile
  buildArgs:
    ALPINE_IMAGE: alpine:3.10
- name: nginx-custom
  tags:
  - "custom-build"
  context: .
  dockerfile: nginx/Dockerfile
```

> **Warning:** don't forget that the Dockerfile path (`.images[].dockerfile`) is relative to the context (`.images[].context`)!
> Also, don't forget the context is relative to which directory you run `bob` command!

## Contributing

For now we are just getting Pull Requests as contribution. There's not much defined yet.

## License

See [LICENSE](LICENSE) for more details.

[⬆ Back to the top](#bob-the-builder)<br>
