# Bully Algorithm Visualization

[![codecov](https://codecov.io/gh/TimTosi/bully-algorithm/branch/master/graph/badge.svg)](https://codecov.io/gh/TimTosi/bully-algorithm)
[![CircleCI](https://circleci.com/gh/TimTosi/bully-algorithm.svg?style=shield)](https://circleci.com/gh/TimTosi/bully-algorithm)
[![Go Report Card](https://goreportcard.com/badge/github.com/timtosi/bully-algorithm)](https://goreportcard.com/report/github.com/timtosi/bully-algorithm)
[![GoDoc](https://godoc.org/github.com/timtosi/bully-algorithm?status.svg)](https://godoc.org/github.com/timtosi/bully-algorithm)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

<img src="assets/intro.gif" alt="intro logo" width="500" />

## Table of Contents
- [About](#what-this-repository-is-about)
- [Bully Algorithm](#what-is-the-bully-algorithm)
- [Quickstart with Docker](#quickstart-with-docker)
- [Quickstart with binaries](#quickstart-without-docker)
- [FAQ](#faq)
- [Support & Feedbacks](#not-good-enough)


## What this repository is about ?

This repository contains source code of an implementation of the bully algorithm
written in Go and a small browser visualization tool.

This has been made for learning purposes about [distributed algorithms](https://en.wikipedia.org/wiki/Distributed_algorithm), Bully algorithm being the simplest leader election algorithm to implement.

Finally, I feel like implementing an algorithm myself helps me to understands it
better and I thought it could be interesting to someone else.


## What is the Bully algorithm ?

The [Bully algorithm](https://en.wikipedia.org/wiki/Bully_algorithm) is one of
the simplest algorithm made to design a coordinator among a set of machines.


## Quickstart

First, go get this repository:
```sh
go get -d github.com/timtosi/bully-algorithm
```

### Quickstart with Docker

> :exclamation: If you don't have [Docker](https://docs.docker.com/install/) and
> [Docker Compose](https://docs.docker.com/compose/) installed, you still can
> execute this program by [compiling the binaries](#quickstart-without-docker). 

This program comes with an already configured [Docker Compose](https://github.com/TimTosi/bully-algorithm/blob/master/deployments/docker-compose.yaml)
file launching five nodes and the browser based user interface.

You can use the `run` target in the provided Makefile to use it easily:

[![asciicast](https://asciinema.org/a/228925.svg)](https://asciinema.org/a/228925)

You can access the visualization through your browser at `localhost:8080`.
If you want to test the cluster behaviour, you can stop and resume some of the
nodes with docker commands.

> :bulb: If you want to update the number of node or change some IDs you will
> have to update the [configuration file](https://github.com/TimTosi/bully-algorithm/blob/master/cmd/bully/conf/bully.conf.yaml#L14-L19)
> and the [Docker Compose file](https://github.com/TimTosi/bully-algorithm/blob/master/deployments/docker-compose.yaml)
> accordingly.


### Quickstart without Docker

First compiles and launch the visualization server:
```sh
cd $GOPATH/src/github.com/timtosi/bully-algorithm/cmd/data-viz
go build && ./data-viz
```

![Visu](assets/run-visu.gif)

Then launch at least two nodes with specifying their ID in argument:
```sh
cd $GOPATH/src/github.com/timtosi/bully-algorithm/cmd/bully
go build && ./bully 0
```

> :bulb: IDs should by default be comprised between 0 to 4 but you should be
> able to update [peer address default configuration](https://github.com/TimTosi/bully-algorithm/blob/master/cmd/bully/conf.go#L23-L27)
> easily.


![Nodes](assets/run-nodes.gif)

You can access the visualization through your browser at `localhost:8080`.


## FAQ

None so far :raised_hands:


## License

Every file provided here is available under the [MIT License](http://opensource.org/licenses/MIT).


## Not Good Enough ?

If you encouter any issue by using what is provided here, please
[let me know](https://github.com/TimTosi/bully-algorithm/issues) ! 
Help me to improve by sending your thoughts to timothee.tosi@gmail.com !
