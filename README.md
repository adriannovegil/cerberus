
# Cerberus

Monitor your URIs from your any place and get notified when your system is down or response time is more than expected.

## Build the project

To build the project, just execute the following command:

```
 $ make clean; make build
```

The command download the dependencies, pass the Golang linter and build the binary file.

If you want to see more options to build, test, etc. the project execute the command:

```
 $ make help
```

## Use with Docker

If you want to use the app inside a Docker container you can build the Docker image using the following command:

```
 $ make docker
```

Once the image is ready, you can execute it with the following command:

```
 $ docker run -ti local/cerberus
```

## Code Conventions

We recommend using [official go code conventions](https://golang.org/doc/effective_go.html).

## Project Layout

We follow the layout [Standard Go Project Layout
](https://github.com/golang-standards/project-layout) for Go application projects.

It's not an official standard defined by the core Go dev team; however, it is a set of common historical and emerging project layout patterns in the Go ecosystem.

Some of these patterns are more popular than others. It also has a number of small enhancements along with several supporting directories common to any large enough real world application.

## Contributing

For a complete guide to contributing to the project, see the [Contribution Guide](CONTRIBUTING.md).

We welcome contributions of any kind including documentation, organization, tutorials, blog posts, bug reports, issues, feature requests, feature implementations, pull requests, answering questions on the forum, helping to manage issues, etc.

The project community and maintainers are very active and helpful, and the project benefits greatly from this activity.

### Reporting Issues

If you believe you have found a defect in the tool or its documentation, use the repository issue tracker to report the problem to the project maintainers.

If you're not sure if it's a bug or not, start by asking in the discussion forum. When reporting the issue, please provide the version.

### Submitting Patches

The project welcomes all contributors and contributions regardless of skill or experience level.

If you are interested in helping with the project, we will help you with your contribution.

We want to create the best possible tool for our development teams and the best contribution experience for our developers, we have a set of guidelines which ensure that all contributions are acceptable.

The guidelines are not intended as a filter or barrier to participation. If you are unfamiliar with the contribution process, the team will help you and teach you how to bring your contribution in accordance with the guidelines.

For a complete guide to contributing, see the [Contribution Guide](CONTRIBUTING.md).

## License

See the [LICENSE](./LICENSE) file
