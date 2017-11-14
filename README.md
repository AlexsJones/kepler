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
```

>>> help
npm sub commands:
	[npm] file: relink an npm package locally<prefix> <string>
	[npm] remove: remove a dep from package.json <string>
	[npm] usage: find usage of a package within submodules <string>
github sub commands:
	[github] pr: pr command palette
		[pr] attach: attach the current issue to a pr <owner> <reponame> <prnumber>
		[pr] create: create a pr <owner> <repo> <base> <head> <title>
	[github] issue: Issue commands
		[issue] create: set the current working issue <owner> <repo> <issuename>
		[issue] set: set the current working issue <issue number>
		[issue] unset: unset the current working issue
		[issue] show: show the current working issue
		[issue] palette: Manipulate the issue palette of working repos
			[palette] add: Add a repository to the palette as part of current working issue by name <name>
			[palette] remove: Remove a repository from the palette as part of the current working issue by name <name>
			[palette] show: Show repositories in the palette as part of the current working issue
			[palette] delete: Delete all repositories in the palette as part of the current working issue
	[github] login: use an access token to login to github
	[github] fetch: fetch remote repos
submodule sub commands:
	[submodule] branch: branch command palette
	[submodule] exec: execute in all submodules <command string>
storage sub commands:
	[storage] clear: clear all data from kepler
	[storage] show: show storage data
palette sub commands:
	[palette] branch: switch branches or create if they don't exist for working issue palette repos <branchname>
		[branch] push: For pushing the local branches to new/existing remotes
		[branch] local: For switching local branches on palette repos
	[palette] show: Show repositories in the palette as part of the current working issue


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
