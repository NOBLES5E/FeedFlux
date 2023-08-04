# FeedFlux

FeedFlux is a light and versatile tool developed in Go, designed to parse feeds such as RSS, Atom, and more. With FeedFlux, these feeds can be transformed into a unified JSON format and streamed directly to your stdout. Rather uniquely, this tool also offers the ability to record and resume your progress as needed.

## Installation

The built binaries can be downloaded on https://github.com/NOBLES5E/FeedFlux/releases. 

FeedFlux also provides an installation script which is particularly useful in scenarios like CI/CD pipelines. Many thanks to GoDownloader for enabling the easy generation of this script.

By default, it installs in the `./bin` directory relative to the working directory:

```sh
$ sh -c "$(curl --location https://raw.githubusercontent.com/NOBLES5E/FeedFlux/main/install.sh)" -- -d
```

You can override the default installation directory using the `-b` parameter. On Linux, common choices are `~/.local/bin` and `~/bin` to install for the current user, or `/usr/local/bin` to install for all users:

```sh
$ sh -c "$(curl --location https://raw.githubusercontent.com/NOBLES5E/FeedFlux/main/install.sh)" -- -d -b ~/.local/bin 
```

This script makes the installation process easier, especially for automated processes such as continuous integration and continuous deployment.

## Usage

To make use of FeedFlux, specify the feeds you wish to parse in the form of arguments. FeedFlux will fetch these feeds, convert them into a unified JSON format, and stream the output to stdout.

Example:

```sh
$ ff https://example.com/rss https://example.com/atom
```

The example above fetches feeds from the specified URLs.

FeedFlux also includes the functionality to record your fetching progress. When interrupted, FeedFlux can document the current state of feed fetching and enable you to resume later.

Example with recorded progress:

```sh
$ ff -r ./progress/ https://example.com/rss
```

In this case, FeedFlux will use the directory `./progress` to store the progress.

## Examples

### Fetch and print to stdout:

To fetch feed(s) and print the formatted output to stdout; use:

```sh
$ ff https://example.com/rss https://example.com/atom
```
### Record progress:

To fetch feed(s) and record the progress in a JSON file for later resumption; use:

```sh
$ ff -r ./progress https://example.com/rss
```
### Resume fetching:

To continue fetching feed(s) from a previously recorded point, use:

```sh
$ ff -r ./progress -c
```
## Contributing

Contributions to this project are welcomed and appreciated.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
