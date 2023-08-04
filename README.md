# FeedFlux

FeedFlux is a light and versatile tool developed in Go, designed to parse feeds such as RSS, Atom, and more. With FeedFlux, these feeds can be transformed into a unified JSON format and streamed directly to your stdout. Rather uniquely, this tool also offers the ability to record and resume your progress as needed.

## Installation

To install FeedFlux, ensure a working Go environment is set up. If one isn't available, follow the installation instructions on the official Go site.

Upon setup of Go, you can download FeedFlux using the `go get` command:

```sh
$ go get -u github.com/NOBLES5E/feedflux
```

## Usage

To make use of FeedFlux, specify the feeds you wish to parse in the form of arguments. FeedFlux will fetch these feeds, convert them into a unified JSON format, and stream the output to stdout.

Example:

```sh
$ feedflux https://example.com/rss https://example.com/atom
```

The example above fetches feeds from the specified URLs.

FeedFlux also includes the functionality to record your fetching progress. When interrupted, FeedFlux can document the current state of feed fetching and enable you to resume later.

Example with recorded progress:

```sh
$ feedflux -r progress.json https://example.com/rss
```

In this case, FeedFlux will use a file `progress.json` to store the progress.

## Examples

### Fetch and print to stdout:

To fetch feed(s) and print the formatted output to stdout; use:

```sh
$ feedflux https://example.com/rss https://example.com/atom
```
### Record progress:

To fetch feed(s) and record the progress in a JSON file for later resumption; use:

```sh
$ feedflux -r progress.json https://example.com/rss
```
### Resume fetching:

To continue fetching feed(s) from a previously recorded point, use:

```sh
$ feedflux -r progress.json -c
```
## Contributing

Contributions to this project are welcomed and appreciated.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
