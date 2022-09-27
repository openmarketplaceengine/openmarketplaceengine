<div id="top"></div>

<h3 align="center">Open Marketplace Engine</h3>

<p align="center">
  A matchmaking system for labor platforms.
</p>

[![Lines Of Code](https://tokei.rs/b1/github/openmarketplaceengine/openmarketplaceengine?category=code)](https://github.com/openmarketplaceengine/openmarketplaceengine)
[![Build & Test](https://github.com/driverscooperative/geosrv/actions/workflows/build-test.yml/badge.svg)](https://github.com/driverscooperative/geosrv/actions/workflows/build-test.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/openmarketplaceengine/openmarketplaceengine.svg)](https://github.com/openmarketplaceengine/openmarketplaceengine)
[![GoReportCard](https://goreportcard.com/badge/github.com/openmarketplaceengine/openmarketplaceengine)](https://goreportcard.com/report/github.com/openmarketplaceengine/openmarketplaceengine)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#project-status">Project Status</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#api">API</a></li>
    <li><a href="#deployment">Deployment</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgements">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About the Project

Open Marketplace Engine (OME) will provide matching and marketplace services to labor platforms such as rideshare services. OME integrates with other components of a gig labor platform to track and store realtime data such as worker location pings and incoming job requests. The service routes jobs to available workers and serves as a source of truth for the job state machine.

OME is designed to do the heavy lifting of a labor platform, including providing geography services for e.g. rider/driver matching. OME is _not_ intended to provide an entire labor platform; features like user models, authentication, and user-facing interfaces are non-goals.

OME is being developed for use in the rideshare platform developed by The Drivers Cooperative. If you are interested in adopting OME in your platform, please [get in touch](mailto:jason@drivers.coop).


## Project Status

As of January 2022, OME has begun active development. Tentative milestones include:

- [Milestone 1: Worker status & pings](https://github.com/orgs/openmarketplaceengine/projects/1/views/1?layout=board)
- Milestone 2: Point-to-point job tracking and assignment
- Milestone 3: Persistence
- Milestone 4: Basic geographic model
- Milestone 5: Basic dispatch


## Getting Started

Dependencies can be launched using `docker-compose`:

```shell
$ docker-compose up
```

You can run the OME service itself with:
```shell
OME_PGDB_ADDR=postgres://postgres:$(whoami)@localhost:5432/$(whoami) go run main.go
```

### Makefile

Run `make` to list all make tasks.

```shell
$ make

echo-env                       Echo environment variables
clean                          Clean
nuke                           Clean with -modcache
build                          Build binary
checkstyle                     Run quick checkstyle (govet + goimports (fail on errors))
test                           Run tests
test-cover                     Run tests with -covermode
test-race                      Run tests with -race
test-bench                     Run tests with -bench
```

### Configuration

_Coming soon_

```shell
export OME_SERVICE_PORT='10000'

export OME_REDIS_STORE_HOST='localhost:6379'
export OME_REDIS_STORE_USERNAME=default
export OME_REDIS_STORE_PASSWORD=secret
```


## API

OME will publish a [gRPC](https://grpc.io/) API for updating worker and job states. The API is under development.


## Deployment

OME will be available as a container image. [Track the issue](https://github.com/driverscooperative/geosrv/issues/4)


## Contributing

Contributions are welcome! Please communicate actively on the [projects board](https://github.com/orgs/openmarketplaceengine/projects?type=beta) before beginning work, as most coordination is currently being done in sidechannels.

_Coming soon: Contributor guidelines, code of conduct, and CLA_


## License

Assigning a license to the OME source is currently a work-in-progress. Inspiration is being taken from the [CoopyLeft License](https://wiki.coopcycle.org/en:license) and the [Business Source License](https://mariadb.com/bsl11/).


## Contact

To get in touch with the OME team, you can:

- [File an issue](https://github.com/driverscooperative/geosrv/issues/new)
- [Email the technical coordinator](mailto:jason@drivers.coop)

As the project picks up we will consider a Discord/Slack setup as well.


## Acknowledgements

OME is being developed and funded by [The Drivers Cooperative](https://drivers.coop), a driver-owned rideshare cooperative in New York City, USA. You can donate to the operations of the cooperative [here](https://ioby.org/project/system-change-rideshare-platform-economy).
