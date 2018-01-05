```
_  _  ____  ____  __    ____  ____
( )/ )( ___)(  _ \(  )  ( ___)(  _ \
 )  (  )__)  )___/ )(__  )__)  )   /
(_)\_)(____)(__)  (____)(____)(_)\_)
```

[![GoDoc](https://godoc.org/github.com/AlexsJones/kepler?status.svg)](https://godoc.org/github.com/AlexsJones/kepler)
[![GitHub issues](https://img.shields.io/github/issues/AlexsJones/kepler.svg)](https://github.com/AlexsJones/kepler/issues)
[![GitHub stars](https://img.shields.io/github/stars/AlexsJones/kepler.svg)](https://github.com/AlexsJones/kepler/stargazers)

Designed to do a few things well...

- work with github issues & pr's.
- let you run multiple submodule commands.
- work with npm modules within submodules.
- works with kubebuilder to deploy kubernetes projects (https://github.com/AlexsJones/kubebuilder.git)
- auto updater built in.
```
>>> help
[kubebuilder]: kubebuilder command palette
	[auth]: Authenticates you into GCP GCR
	[build]: Builds a docker image based off a kepler definitions
	[deploy]: Deploy to a remote kubebuilder cluster
[node]: node command palette
	[remove]: remove a dep from package.json <string>
	[usage]: find usage of a package within submodules <string>
	[view]: View all the node projects that can be found locally
	[local-deps]: Shows all the dependancies found locally
	[install]: Installs all the required vendor code
	[init]: Create the package json for a meta repo
[github]: github command palette
	[team]: team command palette
		[list]: List team membership
		[set]: Set the current team to work with
		[fetch]: Fetch remote team repos
	[pr]: pr command palette
		[attach]: attach the current issue to a pr <owner> <reponame> <prnumber>
		[create]: create a pr <owner> <repo> <base> <head> <title>
	[issue]: Issue commands
		[create]: set the current working issue <owner> <repo> <issuename>
		[set]: set the current working issue <issue number>
		[unset]: unset the current working issue
		[show]: show the current working issue
		[palette]: Manipulate the issue palette of working repos
	[add]: Add a repository to the palette as part of current working issue by name <name>
	[remove]: Remove a repository from the palette as part of the current working issue by name <name>
	[show]: Show repositories in the palette as part of the current working issue
	[delete]: Delete all repositories in the palette as part of the current working issue
	[login]: use an access token to login to github
	[fetch]: fetch remote repos
[submodule]: submodule command palette
	[branch]: branch command palette
	[exec]: execute in all submodules <command string>
[storage]: storage command palette
	[clear]: clear all data from kepler
	[show]: show storage data
[palette]: Issue palette that controls repos can be used from here
	[branch]: switch branches or create if they don't exist for working issue palette repos <branchname>
		[push]: For pushing the local branches to new/existing remotes
		[local]: For switching local branches on palette repos
	[show]: Show repositories in the palette as part of the current working issue
[docker]: docker command palette
	[build]: Builds a project in standalone from the defined Dockerfile
```
[![asciicast](https://asciinema.org/a/uccLCSINhgn48JBMFMEDNLCZg.png)](https://asciinema.org/a/uccLCSINhgn48JBMFMEDNLCZg)


### Working with github issues and selected repos with `palette`

[![asciicast](https://asciinema.org/a/QRO4nWiZycVLi8jZb9TFnUbau.png)](https://asciinema.org/a/QRO4nWiZycVLi8jZb9TFnUbau)



## Installation

Either:

```
go get github.com/AlexsJones/kepler
```

or fetch source and run
```
godep restore
```

## Non-CLI mode

To run in a CI or as a one time command

```
kepler unattended usage yar
```

## Development

[![Maintainability](https://api.codeclimate.com/v1/badges/31068bf57e3db317466b/maintainability)](https://codeclimate.com/github/AlexsJones/kepler/maintainability)

As you can see maintainability is important and constant refactoring and improvement shall be made.
