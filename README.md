# Open Marketplace Engine

## How to start

```shell
docker-compose up
```
## Makefile

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

# Config

Environment variable override example.
```shell
export OME_SERVICE_PORT='10000'

export OME_REDIS_STORE_HOST='localhost:6379'
export OME_REDIS_STORE_USERNAME=default
export OME_REDIS_STORE_PASSWORD=secret
```


