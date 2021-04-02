package testutil

import "github.com/aws-controllers-k8s/dev-tools/pkg/config"

// Returns a new config.Config object used for testing purposes.
func NewConfig(services ...string) *config.Config {
	return &config.Config{
		Repositories: config.RepositoriesConfig{
			Core: []string{
				"runtime",
				"code-generator",
			},
			Services: services,
		},
		Github: config.GithubConfig{
			ForkPrefix: "ack-",
			Username:   "SrinivasaRamanujan",
		},
	}
}
