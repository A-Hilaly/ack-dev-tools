package repository

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	git "gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	ackdevgit "github.com/aws-controllers-k8s/dev-tools/pkg/git"
	"github.com/aws-controllers-k8s/dev-tools/pkg/github"
	"github.com/aws-controllers-k8s/dev-tools/pkg/util"
)

const (
	originRemoteName   = "origin"
	upstreamRemoteName = "upstream"
)

var (
	ErrUnconfiguredRepository   error = errors.New("unconfigured repository")
	ErrRepositoryDoesntExist    error = errors.New("unknown doesnt exist")
	ErrRepositoryAlreadyExist   error = errors.New("repository already exist")
	ErrMissingGithubCredentials error = errors.New("missing github credentials")
)

// NewManager create a new manager.
func NewManager(cfg *config.Config, opts ...Option) (*Manager, error) {
	ghc := github.NewClient(cfg.Github.Token)
	m := &Manager{
		cfg: cfg,
		ghc: ghc,
		git: ackdevgit.New(nil),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m, nil
}

// Manager is reponsible of managing local ACK local repositories and
// github forks.
type Manager struct {
	log       *logrus.Logger
	cfg       *config.Config
	ghc       github.RepositoryService
	git       ackdevgit.OpenCloner
	repoCache []*Repository
}

// LoadRepository loads information about a single repository
func (m *Manager) LoadRepository(name string, t RepositoryType) (*Repository, error) {
	// check repo cache
	repo, err := m.GetRepository(name)
	if err == nil {
		return repo, nil
	}

	// fail if repository doesn't exist in the manager configuration
	switch t {
	case RepositoryTypeCore:
		if !util.InStrings(name, m.cfg.Repositories.Core) {
			return nil, ErrUnconfiguredRepository
		}
	case RepositoryTypeController:
		if !util.InStrings(name, m.cfg.Repositories.Services) {
			return nil, ErrUnconfiguredRepository
		}
	}

	// controller repositories should always have a '-controller' suffix
	if t == RepositoryTypeController {
		name = fmt.Sprintf("%s-controller", name)
	}

	// set expected fork name
	forkName := name
	if m.cfg.Github.ForkPrefix != "" {
		forkName = fmt.Sprintf("%s%s", m.cfg.Github.ForkPrefix, name)
	}

	fullPath := filepath.Join(m.cfg.RootDirectory, name)
	var gitRepo *git.Repository
	var gitHead string

	gitRepo, err = m.git.Open(fullPath)
	if err != nil && err != git.ErrRepositoryNotExists {
		return nil, err
	} else if err == nil {
		head, err := gitRepo.Head()
		if err != nil {
			return nil, err
		}
		gitHead = head.Name().Short()
	}

	return &Repository{
		Name:      name,
		Type:      t,
		gitRepo:   gitRepo,
		GitHead:   gitHead,
		FullPath:  fullPath,
		RemoteURL: remoteURL(m.cfg.Github.Username, forkName),
	}, nil
}

// LoadAll parses the configuration and loads informations about local
// repositories if they are found.
func (m *Manager) LoadAll() error {
	// collect repositories from config
	for _, coreRepo := range m.cfg.Repositories.Core {
		repo, err := m.LoadRepository(coreRepo, RepositoryTypeCore)
		if err != nil {
			return err
		}
		m.repoCache = append(m.repoCache, repo)
	}
	for _, serviceName := range m.cfg.Repositories.Services {
		repo, err := m.LoadRepository(serviceName, RepositoryTypeController)
		if err != nil {
			return err
		}
		m.repoCache = append(m.repoCache, repo)
	}
	return nil
}

// GetRepository return a known repository
func (m *Manager) GetRepository(repoName string) (*Repository, error) {
	for _, repo := range m.repoCache {
		if repo.Name == repoName {
			return repo, nil
		}
	}
	return nil, ErrRepositoryDoesntExist
}

// List lists all the cached repositories. Alias of ListAnd().
func (m *Manager) List(filters ...Filter) []*Repository {
	return m.ListAnd(filters...)
}

// List lists all the cached repositories
func (m *Manager) ListAnd(filters ...Filter) []*Repository {
	repos := []*Repository{}
mainLoop:
	for _, repo := range m.repoCache {
		for _, filter := range filters {
			if !filter(repo) {
				continue mainLoop
			}
		}
		repos = append(repos, repo)
	}
	return repos
}

// List lists all the cache repositories
func (m *Manager) ListOr(filters ...Filter) []*Repository {
	repos := []*Repository{}
mainLoop:
	for _, repo := range m.repoCache {
		for _, filter := range filters {
			if filter(repo) {
				repos = append(repos, repo)
				continue mainLoop
			}
		}
	}
	return repos
}

// Clone clones a known repository to the root directory
func (m *Manager) clone(repoName string) error {
	repo, err := m.GetRepository(repoName)
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %v", repoName, err)
	}
	if repo.gitRepo != nil {
		return ErrRepositoryAlreadyExist
	}

	// clone fork repository
	err = m.git.Clone(context.TODO(), repo.ForkURL, repo.FullPath)
	if errors.Is(err, transport.ErrAuthenticationRequired) {
		return ErrMissingGithubCredentials
	}
	if err != nil {
		return fmt.Errorf("cannot clone repository %s: %v", repoName, err)
	}

	// open git repository
	gitRepo, err := m.git.Open(repo.FullPath)
	if err != nil {
		return err
	}

	repo.gitRepo = gitRepo
	// Add upstream remote
	gitRepo.CreateRemote(&gitconfig.RemoteConfig{
		Name: upstreamRemoteName,
		URLs: []string{repo.RemoteURL},
	})

	return nil
}

// ensureFork ensures that your github account have a fork for a given
// ACK project. It will also rename the project if it's not following the
// standard: $ackprefix-$projectname
func (m *Manager) ensureFork(repo *Repository) error {
	m.log.SetLevel(logrus.DebugLevel)

	// TODO(hilalymh) maybe we should propagate the context from the cobra commands
	ctx := context.TODO()
	expectedForkName := fmt.Sprintf("%s%s", m.cfg.Github.ForkPrefix, repo.Name)

	fork, err := m.ghc.GetUserRepositoryFork(ctx, repo.Name)
	if err == nil {
		if *fork.Name != expectedForkName {
			err = m.ghc.RenameRepository(ctx, m.cfg.Github.Username, *fork.Name, expectedForkName)
			if err != nil {
				return err
			}
			repo.ForkName = expectedForkName
		}
	} else if err == github.ErrorForkNotFound {
		err = m.ghc.ForkRepository(ctx, repo.Name)
		if err != nil {
			return err
		}

		time.Sleep(1 * time.Second)

		err = m.ghc.RenameRepository(ctx, m.cfg.Github.Username, repo.Name, expectedForkName)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func (m *Manager) ensureClone(repo *Repository) error {
	err := m.clone(repo.Name)
	if err != nil && err != ErrRepositoryAlreadyExist {
		return err
	}
	return nil
}

// EnsureAll ensures one repository.
func (m *Manager) EnsureRepository(name string) error {
	repo, err := m.GetRepository(name)
	if err != nil {
		return err
	}

	err = m.ensureFork(repo)
	if err != nil {
		return err
	}

	err = m.ensureClone(repo)
	if err != nil {
		return err
	}

	return nil
}

// EnsureAll ensures all cached repositories.
func (m *Manager) EnsureAll() error {
	for _, repo := range m.repoCache {
		err := m.ensureFork(repo)
		if err != nil {
			return err
		}

		err = m.ensureClone(repo)
		if err != nil {
			return err
		}
	}
	return nil
}
