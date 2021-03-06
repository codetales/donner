# Donner
[![Go Report Card](https://goreportcard.com/badge/github.com/codetales/donner)](https://goreportcard.com/report/github.com/codetales/donner)
[![Build Status](https://travis-ci.com/codetales/donner.svg?branch=master)](https://travis-ci.com/codetales/donner)

:boom: Donner is a generic command wrapper. It let's you define strategies to wrap commands in things like `docker-compose exec` or `docker container run`.
This is can come in very handy when developing applications in containers.  
Donner allows defining a wrapping strategy on a per command basis. So you don't have to worry which service to use or whether you should use `docker-compose exec` or `docker-compose run` when executing a command.

## Installation

`brew install codetales/tap/donner` or download latest binary from the [releases](https://github.com/codetales/donner/releases) page.

## Examples

Example config for a ruby project
```yaml
strategies:
  run:
    handler: docker_compose_run
    service: app
    remove: true
  exec:
    handler: docker_compose_exec
    service: app
  exec_postgres:
    handler: docker_compose_exec
    service: pg
  run_with_docker:
    handler: docker_run
    image: alpine:latest
    volumes:
      - "./:/usr/src/app"

default_strategy: run       # use this strategy for all undefined commands

commands:
  rails: exec
  rspec: exec
  ruby: run
  bundle: run
  irb: run
  rake:                     # override strategy definition
    strategy: run
    remove: false
  psql: exec_postgres
```

## Running stuff
```
donner run ruby -v                # Executes `docker-compose run --rm app ruby -v`
donner run rails c                # Executes `docker-compose exec app rails c`
donner run psql                   # Executes `docker-compose exec pg psql`
donner run ./bin/rspec spec/      # Executes `docker-compose exec app ./bin/rspec spec/` - Donner strips of the path when checking an executable
donner run rake my:task           # Executes `docker-compose run app rake my:task`
donner run irb                    # Executes `docker-compose run --rm app irb`
donner run other-command          # Executes `docker-compose run --rm app other-command`
```

## Setting aliases
```
$ donner aliases

alias rails='donner run rails'
alias rspec='donner run rspec'
alias rake='donner run rake'
alias psql='donner run psql'
alias ruby='donner run ruby'
alias bundle='donner run bundle'
alias irb='donner run irb'

# copy and paste the output into your terminal or run
#   eval $(donner aliases)
```

## Fallbacks and strict checking
To make aliases work nicely system wide, Donner supports a `--fallback` and `--strict` flag.

* When executing `donner run --strict some-command` and the provided command is not defined under `commands` or `aliases`, Donner will fail executing the command.
* When executing `donner run --fallback` with no `.donner.yml` in the current directory, Donner executes the command as is.
* When executing `donner run --strict --fallback` and the provided command is not defined, Donner executes the command as is.


This is also supported with aliases:
```
$ donner aliases --fallback --strict

alias rails='donner run --fallback --strict rails'
alias rspec='donner run --fallback --strict rspec'
alias rake='donner run --fallback --strict rake'
alias psql='donner run --fallback --strict psql'
alias ruby='donner run --fallback --strict ruby'
alias bundle='donner run --fallback --strict bundle'
alias irb='donner run --fallback --strict irb'

# copy and paste the output into your terminal or run
#   eval $(donner aliases --fallback --strict)
```

## TODO
* Ensure useful error messages and test failure cases
* Add missing flags to the handlers (volumes, ...)
* Add documentation
* Setup bash/zsh completion
