package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	"github.com/aws-controllers-k8s/dev-tools/pkg/repository"
)

func init() {}

var ensureRepositoriesCmd = &cobra.Command{
	Use:     "repo",
	Aliases: []string{"repo", "repos", "repositories"},
	RunE:    ensureAll,
	Args:    cobra.NoArgs,
	Short:   "Ensure repositories are forked and cloned locally",
	Long: `Ensure repositories are forked and cloned locally.
This command will also rename forks and local clones if you have configured 
a repo name prefix.`,
}

func ensureAll(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(ackConfigPath)
	if err != nil {
		return err
	}

	repoManager, err := repository.NewManager(cfg)
	if err != nil {
		return err
	}

	err = repoManager.LoadAll()
	if err != nil {
		return err
	}

	ctx := context.TODO()
	err = repoManager.EnsureAll(ctx)
	if err != nil {
		return err
	}

	return nil
}
