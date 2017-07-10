```
_  _  ____  ____  __    ____  ____
( )/ )( ___)(  _ \(  )  ( ___)(  _ \
 )  (  )__)  )___/ )(__  )__)  )   /
(_)\_)(____)(__)  (____)(____)(_)\_)
```

Designed to do a few things well...
- work with github issues & pr's.
- let you run multiple submodule commands.
- work with npm modules within submodules.
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
	[github] issue: Issue command palette
		[issue] create: set the current working issue <owner> <repo> <issuename>
		[issue] set: set the current working issue <issue url>
		[issue] unset: unset the current working issue
		[issue] show: show the current working issue
	[github] login: use an access token to login to github
submodule sub commands:
	[submodule] branch: branch command palette
	[submodule] exec: execute in all submodules <command string>
storage sub commands:
	[storage] clear: clear all data from kepler
	[storage] show: show storage data


```
[![asciicast](https://asciinema.org/a/uccLCSINhgn48JBMFMEDNLCZg.png)](https://asciinema.org/a/uccLCSINhgn48JBMFMEDNLCZg)


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
