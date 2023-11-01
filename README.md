# Go library for the GoCtl CLI

Modules from this library will obeyGoCtl CLI conventions by default:

- [`repository.Current()`](https://pkg.go.dev/github.com/khulnasoft-lab/go-goctl/v2/pkg/repository#current) respects the value of the `GOCTL_REPO` environment variable and reads from git remote configuration as fallback.

- GitHub API requests will be authenticated using the same mechanism as `gh`, i.e. using the values of `GOCTL_TOKEN` and `GOCTL_HOST` environment variables and falling back to the user's stored OAuth token.

- [Terminal capabilities](https://pkg.go.dev/github.com/khulnasoft-lab/go-goctl/v2/pkg/term) are determined by taking environment variables `GOCTL_FORCE_TTY`, `NO_COLOR`, `CLICOLOR`, etc. into account.

- Generating [table](https://pkg.go.dev/github.com/khulnasoft-lab/go-goctl/v2/pkg/tableprinter) or [Go template](https://pkg.go.dev/github.com/khulnasoft-lab/go-goctl/pkg/template) output uses the same engine as gh.

- The [`browser`](https://pkg.go.dev/github.com/khulnasoft-lab/go-goctl/v2/pkg/browser) module activates the user's preferred web browser.

## Usage

See the full `go-goctl`  [reference documentation](https://pkg.go.dev/github.com/khulnasoft-lab/go-goctl/v2) for more information

```golang
package main

import (
	"fmt"
	"log"
	"github.com/khulnasoft-lab/go-goctl/v2"
	"github.com/khulnasoft-lab/go-goctl/v2/pkg/api"
)

func main() {
	// These examples assume `gh` is installed and has been authenticated.

	// Shell out to a goctl command and read its output.
	issueList, _, err := gh.Exec("issue", "list", "--repo", "khulnasoft-lab/goctl", "--limit", "5")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(issueList.String())

	// Use an API client to retrieve repository tags.
	client, err := api.DefaultRESTClient()
	if err != nil {
		log.Fatal(err)
	}
	response := []struct{
		Name string
	}{}
	err = client.Get("repos/khulnasoft-lab/goctl/tags", &response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}
```

See [examples][] for more demonstrations of usage.

## Contributing

If anything feels off, or if you feel that some functionality is missing, please check out our [contributing docs][contributing]. There you will find instructions for sharing your feedback and for submitting pull requests to the project. Thank you!

[extensions]: https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions
[examples]: ./example_goctl_test.go
[contributing]: ./.github/CONTRIBUTING.md
