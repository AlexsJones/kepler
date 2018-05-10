package main

import (
	"errors"

	"github.com/AlexsJones/kepler/commands/storage"
	"github.com/DATA-DOG/godog"
)

var (
	storageRef *storage.Storage
)

func iGetAnInstanceOfStorage() error {
	storageRef = storage.GetInstance()
	return nil
}

func iAmAbleToValidateThatTheInstanceIsAnInitialisedObject() error {
	if storageRef == nil {
		return errors.New("Uninitialised storage")
	}
	return nil
}

func StorageFeatureContext(s *godog.Suite) {
	s.Step(`^I get an instance of storage$`, iGetAnInstanceOfStorage)
	s.Step(`^I am able to validate that the instance is an initialised object$`, iAmAbleToValidateThatTheInstanceIsAnInitialisedObject)

}
func iWishToStoreAnAccessTokenWithValueOf(arg1 string) error {
	storageRef.Github.AccessToken = arg1
	return storageRef.Save()
}

func itShouldSaveCorrectly(args1 string) error {
	if storage.GetInstance().Github.AccessToken != args1 {
		return errors.New("Stored value does not match")
	}
	return nil
}

func StorageWithDataFeatureContext(s *godog.Suite) {
	s.Step(`^I wish to store an Access Token with value of "([^"]*)"$`, iWishToStoreAnAccessTokenWithValueOf)
	s.Step(`^it should save correctly with value of "([^"]*)"$`, itShouldSaveCorrectly)
	s.BeforeScenario(func(interface{}) {
		storageRef = storage.GetInstance()
	})
}
