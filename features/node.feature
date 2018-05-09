Feature: Node module validation
  In order to ensure node modules are correctly being manipulated
  As a DevOps engineer
  I need to test several scenarios where node modules are used

  Scenario: Validating a submodule package
  Given I have a generated test submodule with node package.json
  When I attempt to update package contents
  Then I am able to validate the changes within the package have been made
