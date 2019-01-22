# Gomono

> Note: This project was build in a few hours to prove that it can be used in a few specific use
> cases. It may contain bugs and doesn't have a lot of features. It's not
> guaranteed to be stable and it's API may change in the future. If you want
> to use it anyway, you probably want to to fork it.

Gomono is a tool for analyzing changes in a git repository of golang source code and
perform actions like building or testing only the code, which was changed
between two certain commits.

Gomono reads go packages, creates a dependency graph and uses the version history, to find out whether a file in a package or a dependency (either direct or transitive) of a package changed.
Once it is known, which directories contain changes, it's possible to execute a
certain action. Currently only the go build system and make can be invoked.

This project can be used in monorepositories and CI-Pipelines, to improve build
speed, because every build only includes exactly the changes as opposed to the
whole repository, which can grow quite large.

## Usage

Here's the basic usage of gomono:

```shell
gomono [cmd] [flags] [packages]
```

### commands

Currently there's only one command provided, the `build` command. `build` takes
a set of (main) go-packages and reads their dependency graph from the source code.
After the graph is build, there are a few options, which can be set via flags:

* `-action`: action is required and specifies, how to build the changed code.
  For more, see [actions](#actions)
* `-all`: if set, gomono ignores the version history and just uses includes all
  files and directories in the build.
* `-builder`: once an action like for example `test` has been selected, a builde
  specifies the implementation, which is used to actually execute the action.
  This can be something like go test or a Makefile to setup more complex tasks.
  See [builders](#builders)
* `-from`: set the first commit, defaults to `HEAD~1`. This is passed to `git
  diff`, thus it actually doesn't make sense, to refer to a range of commits
  beeing compared, but rather two endpoints (see manpage for `git-diff`).
* `-print`: only prints the dependency graph in dot (see graphviz) format and
  exit (this may be moved to a separated command in the future).
* `-to`: counterpart of `-from`, selects the second commit to be included.

### actions

Actions define a set of steps to be executed, different actions can include
different sets of source code files. For example it may be useful to execute a
test action on all directories, but a deploy action may not be useful for
directories, which contain only libraries. Currently implemented actions are
`deploy`, `build`, bot include only main packages and `test`, includes all go
source files.

### builders

Builders can be used to choose between different implementations of an action.
Tests may be easily executed with `go test`, but deployments can be quite
complex. Available builders are `golang`, to invoke go tools, and `makefile` to
describe more complicated setups. Builders can execute actions.

#### makefile

Makefiles can be used in the same directory as main packages, to implement
`build`, `test` and `deploy` actions. A Makefile has implement a target for the
chosen action.

#### golang

The Golang builder currently only supports the `build` and `test` actions.

## External dependencies

External dependencies can make things complicated. The easiest way probably is to
vendor everything and treat it like usual code controlled by your version control
system. However, as we're using dep to manage our dependencies and it was a 
requirement to keep vendor directories out of source control, gomono tries to
read Gopkg.lock and analyze changes in the used versions, to determine, whether a
project needs to be rebuild.

## Examples

Here are two examples, how gomono may be used in a CI-Pipeline using Gitlab ci,
where `CI_COMMIT_BEFORE_SHA` references the state before the last push and
`GOMONO_ALL` is a custom environment variable, which can be set to true, if a
pipeline ist triggered manually, to rebuild everything:

```shell
gomono build -action=deploy -all=$GOMONO_ALL -from=$CI_COMMIT_BEFORE_SHA -builder=makefile ./...

gomono build -action=test -all=$GOMONO_ALL -from=$CI_COMMIT_BEFORE_SHA -builder=golang ./...
```

The first call will invoke `make deploy` in all directories, which contain a
main project and have changed between the previous commit and the newly pushed
commits. The second will run `go test` for all changed packages.
