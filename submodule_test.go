package main

import (
	"errors"

	"github.com/AlexsJones/kepler/commands/submodules"
	"github.com/DATA-DOG/godog"
)

func SubmoduleFeatureContext(s *godog.Suite) {

	var output error

	s.Step(`^I run a submodule command locally within kepler$`, func() error {

		output = submodules.CommandSubmodules("ll")
		return nil
	})
	s.Step(`^I expect an error code$`, func() error {
		if output == nil {
			return errors.New("Did not correctly detect non-submodule path")
		}
		return nil
	})
}
