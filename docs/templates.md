## Templates

[Example](https://github.com/RedSailTechnologies/boatswain/blob/main/docs/example.yaml)

Templates are used in Boatswain to merge configurations together such that deliveries can be kept dry and easy to reuse. The high level objects are as follows (details about each below):
* Deployment
* Test
* Trigger
* Strategy

At a high level the user configures deployments to utilize helm to install/upgrade applications, tests to ensure the application's validity before moving on to the next steps of a strategy, triggers to start a delivery's strategy, and strategies to specify how the deployments and tests are run once triggered.
