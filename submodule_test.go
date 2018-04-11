package main

import "github.com/DATA-DOG/godog"

func iRunASubmoduleCommandLocallyWithinKepler() error {
	return godog.ErrPending
}

func iExpectAnErrorCode() error {
	return godog.ErrPending
}

func SubmoduleFeatureContext(s *godog.Suite) {
	s.Step(`^I run a submodule command locally within kepler$`, iRunASubmoduleCommandLocallyWithinKepler)
	s.Step(`^I expect an error code$`, iExpectAnErrorCode)
}
