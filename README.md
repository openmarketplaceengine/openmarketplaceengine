# Open Marketplace Engine

Open Marketplace Engine is matching and marketplace service, a component of scalable, efficient and robust labor platform.

## Features
Accepts realtime data related to incoming consumer requests and available laborers, including their realtime geographic locations. 
This data is used to dispatch requests to laborers.

## Quickstart

```shell
docker-compose up
```
Run `make` command to list all make tasks.
```shell
make

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

[config/config.yml](config/config.yml) file values can be overridden by environment variables.
```shell
export OME_SERVICE_PORT='10000'

export OME_REDIS_STORE_HOST='localhost:6379'
export OME_REDIS_STORE_USERNAME=default
export OME_REDIS_STORE_PASSWORD=secret
```


