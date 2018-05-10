package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"

	"github.com/AlexsJones/kepler/commands/node"
	"github.com/DATA-DOG/godog"
)

const (
	p = `{
  "name": "api",
  "version": "3.1.15",
  "description": "Common config",
  "main": "index.js",
  "scripts": {
    "test": "mocha test/index.js"
  },
  "author": "Test <dev@test.com>",
  "private": true,
  "license": "ISC",
  "dependencies": {
    "confidence": "^1.4.2",
    "debug": "^2.3.3",
    "joi": "^10.6.0",
    "lodash": "^4.17.4",
    "shortid": "^2.2.8",
    "swig": "^1.4.2",
    "uuid": "^2.0.2"
  }
}`
)

func iHaveAGeneratedTestSubmoduleWithNodePackagejson() error {
	os.Mkdir("test-submodule-repo", 0755)
	exec.Command("bash", "-c", "cd test-submodule-repo && git init").CombinedOutput()
	_, err := os.Stat("test-submodule-repo")
	if err != nil {
		return err
	}
	exec.Command("bash", "-c", "git submodule add ./test-submodule-repo").CombinedOutput()

	f, err := os.Create("./test-submodule-repo/package.json")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	io.WriteString(w, p)
	w.Flush()
	return nil
}

func iAttemptToUpdatePackageContents() error {

	_, err := node.LocalNodeModules()
	if err != nil {
		return err
	}

	return nil
}

func iAmAbleToValidateTheChangesWithinThePackageHaveBeenMade() error {

	os.RemoveAll("test-submodule-repo")
	os.Remove(".gitmodules")
	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have a generated test submodule with node package\.json$`, iHaveAGeneratedTestSubmoduleWithNodePackagejson)
	s.Step(`^I attempt to update package contents$`, iAttemptToUpdatePackageContents)
	s.Step(`^I am able to validate the changes within the package have been made$`, iAmAbleToValidateTheChangesWithinThePackageHaveBeenMade)

}
