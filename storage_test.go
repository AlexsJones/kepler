package main

import "github.com/DATA-DOG/godog"

func iWantToTestTheStorageCommandsRunAsExpected() error {
	return godog.ErrPending
}

func iAmAmAbleToValidateThisByInitialization() error {
	return godog.ErrPending
}

func StorageFeatureContext(s *godog.Suite) {
	s.Step(`^I want to test the storage commands run as expected$`, iWantToTestTheStorageCommandsRunAsExpected)
	s.Step(`^I am am able to validate this by initialization$`, iAmAmAbleToValidateThisByInitialization)
}
