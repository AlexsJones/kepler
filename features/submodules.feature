Feature: Submodule command validation
  In order to ensure Kepler submodules work correctly
  As a DevOps engineer
  I need to test within the kepler repo they don't fire off false positives


  Scenario: Test submodule parsing does not trigger
    Given I run a submodule command locally within kepler
    Then I expect an error code
