# Bob, the builder

> Bob is an automated multiple Docker images builder. It works by defining a simple YAML file and running its commands.

### Under development

This project is under heavy development.

Tasks:

- [x] Build images
- [ ] Push images

## Installing

### From source

To install the project from its source code you need to clone this repository and have `golang` installed.

```
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
  bob build --file bobber.yaml

Usage:
  bob [command]

Available Commands:
  build       Builds Docker images
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -f, --file string   yaml configuration file (default "bob.yaml")
  -h, --help          help for bob

Use "bob [command] --help" for more information about a command.
```

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

## Contributing

For now we are just getting Pull Requests as contribution. There's not much defined yet.

## License

See [LICENSE](LICENSE.md) for more details.

[⬆ Back to the top](#bob-the-builder)<br>
