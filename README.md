```
_  _  ____  ____  __    ____  ____
( )/ )( ___)(  _ \(  )  ( ___)(  _ \
 )  (  )__)  )___/ )(__  )__)  )   /
(_)\_)(____)(__)  (____)(____)(_)\_)
```

Designed to do a few things well...

```
npm sub commands:
	[npm] file: relink an npm package locally
	[npm] remove: remove a file from package.json
	[npm] usage: find usage of a package within submodules
github sub commands:
	[github] issue: Issue command palette
		[issue] set: set the current working issue <reponame> <issuenumber>
		[issue] unset: unset the current working issue
		[issue] pr: update a pr with issue information
	[github] login: use an access token to login to github
submodule sub commands:
	[submodule] exec: execute in all submodules
storage sub commands:
	[storage] clear: clear all data from kepler

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
