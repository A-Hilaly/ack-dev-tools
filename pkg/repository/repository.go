package repository

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
)

// NewRepository returns a pointer to a new repository.
func NewRepository(name string, repoType RepositoryType) *Repository {
	if repoType == RepositoryTypeController {
		name = fmt.Sprintf("%s-controller", name)
	}
	return &Repository{
		Name: name,
		Type: repoType,
	}
}

// Repository represents an ACK project repository.
type Repository struct {
	gitRepo *git.Repository

	// Name of the ACK upstream repo
	Name string
	// Expected local full path
	FullPath string
	// Git HEAD commit or current branch
	GitHead string
	// Current fork name. Might be different than the parent repo.
	ForkName string
	// User fork URL.
	ForkURL string
	// Repository Type
	Type RepositoryType
	// RemoteURL
	RemoteURL string
}

func remoteURL(owner, name string) string {
	return fmt.Sprintf("git@github.com:%s/%s.git", owner, name)
}
