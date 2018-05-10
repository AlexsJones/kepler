Feature: Basic storage validation
  In order to ensure Kepler works
  As a DevOps engineer
  I need to test storage save/load features work

  Scenario: Storage
  Given I get an instance of storage
  Then I am able to validate that the instance is an initialised object

  Scenario: Storage of Data
  Given I wish to store an Access Token with value of "test-1234"
  Then it should save correctly with value of "test-1234"
