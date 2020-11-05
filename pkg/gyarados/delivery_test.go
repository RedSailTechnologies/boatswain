package gyarados

import (
	"testing"

	pb "github.com/redsailtechnologies/boatswain/rpc/gyarados"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryValidateInvalidYAML(t *testing.T) {
	inputs := []string{
		// name is required
		`
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// version is required
		`
name: mycooldelivery
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// application name/project are required
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// at least one cluster is required
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters: []
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// at least one deployment or test is required
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments: []
tests: []
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// all deployments/tests are valid
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: aname
    helm:
      chart: mycoolchart
tests: []
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
tests:
  - name: aname
    helm:
      chart: mycoolchart
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// at least one trigger is required
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments: []
tests:
  - template: this-is-gonna-break
triggers: []
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments: []
tests:
  - template: this-is-gonna-break
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// all triggers are valid
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - template: this-ones-valid
  - name: my-trigger
    template: i-shouldnt-specify-both
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// at least one strategy step is required
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy: []
`,
		// all strategy steps are valid
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    template: cant-do-this
`,
		// name and all of the above or template and arguments
		`
name: mycooldelivery
arguments: |
  this: fails
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		`
name: mycooldelivery
template: uh-oh-cant-do-both
`,
		`
template: my-template
version: ${{ .Params.version }}
`,
		`
template: my-template
application:
  name: kube-app-label
`,
		`
template: my-template
application:
  project: cool-proj
`,
		`
template: my-template
clusters:
  - dev
  - prod
`,
		`
template: my-template
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
`,
		`
template: my-template
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
`,
		`
template: my-template
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
`,
		`
template: my-template
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
	}

	for _, input := range inputs {
		sut := &Delivery{&pb.Delivery{}}
		YAML(sut, input)
		assert.NotNil(t, sut.Validate(), input)
	}
}

func TestDeliveryValidateValidYAML(t *testing.T) {
	inputs := []string{
		// valid delivery
		`
name: mycooldelivery
version: ${{ .Params.version }}
application:
  name: kube-app-label
  project: cool-proj
clusters:
  - dev
  - prod
deployments:
  - name: my-deployment
    helm:
      chart: achart
      repo: arepo
      version: ${{ .Params.version }}
tests:
  - name: my-deployment-tests
    docker:
      image: docker.io/myrepo/image
      tag: ${{ .Params.version }}
triggers:
  - name: ci-trigger
    web:
      name: mycoolwebtrigger
strategy:
  - name: deploy-the-app
    displayName: Deploy the Application
    always:
      - deployment: my-deployment
`,
		// valid template
		`
template: mytemplate
arguments: |
  woah: look
  at: all
  of:
    these: arguments
`,
	}

	for _, input := range inputs {
		sut := &Delivery{&pb.Delivery{}}
		YAML(sut, input)
		assert.Nil(t, sut.Validate(), input)
	}
}
