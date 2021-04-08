package repository

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-git.v4"

	"github.com/aws-controllers-k8s/dev-tools/pkg/config"
	ackdevgit "github.com/aws-controllers-k8s/dev-tools/pkg/git"
	"github.com/aws-controllers-k8s/dev-tools/pkg/github"
	"github.com/aws-controllers-k8s/dev-tools/pkg/testutil"

	"github.com/aws-controllers-k8s/dev-tools/mocks"
)

func TestManager_LoadRepository(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	testRepo, err := testutil.NewInMemoryGitRepository()
	require.NoError(err)

	fakeGit := &mocks.OpenCloner{}
	{
		fakeGit.On("Open", "runtime").Return(testRepo, nil)
		fakeGit.On("Open", "s3-controller").Return(nil, git.ErrRepositoryNotExists)
		fakeGit.On("Open", "sqs-controller").Return(nil, ErrUnconfiguredRepository)
	}

	type fields struct {
		cfg       *config.Config
		git       ackdevgit.OpenCloner
		repoCache []*Repository
	}
	type args struct {
		name string
		t    RepositoryType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Repository
		wantErr bool
	}{
		{
			name: "repository exists",
			fields: fields{
				cfg: testutil.NewConfig(),
				git: fakeGit,
			},
			args: args{
				name: "runtime",
				t:    RepositoryTypeCore,
			},
			wantErr: false,
		},
		{
			name: "repository doesn't exists",
			fields: fields{
				cfg: testutil.NewConfig(),
				git: fakeGit,
			},
			args: args{
				name: "s3",
				t:    RepositoryTypeController,
			},
			wantErr: true,
		},
		{
			name: "unconfigured repository",
			fields: fields{
				cfg: testutil.NewConfig(),
				git: fakeGit,
			},
			args: args{
				name: "sqs",
				t:    RepositoryTypeController,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				cfg:       tt.fields.cfg,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			_, err := m.LoadRepository(tt.args.name, tt.args.t)
			if (err != nil) != tt.wantErr {
				assert.Fail(fmt.Sprintf("Manager.LoadRepository() error = %v, wantErr %v", err, tt.wantErr))
			}
		})
	}
}

func TestManager_LoadAll(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	testRepo, err := testutil.NewInMemoryGitRepository()
	require.NoError(err)

	fakeGit := &mocks.OpenCloner{}
	fakeGit.On("Open", "runtime").Return(testRepo, nil)
	fakeGit.On("Open", "code-generator").Return(testRepo, nil)
	fakeGit.On("Open", "s3-controller").Return(nil, git.ErrRepositoryNotExists)
	fakeGit.On("Open", "ecr-controller").Return(nil, bytes.ErrTooLarge)

	type fields struct {
		cfg       *config.Config
		git       ackdevgit.OpenCloner
		repoCache []*Repository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all repositories exists",
			fields: fields{
				cfg: testutil.NewConfig(),
				git: fakeGit,
			},
			wantErr: false,
		},
		{
			name: "repository not found",
			fields: fields{
				cfg: testutil.NewConfig("s3"),
				git: fakeGit,
			},
			wantErr: false,
		},
		{
			name: "unexpected repository error",
			fields: fields{
				cfg: testutil.NewConfig("ecr"),
				git: fakeGit,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				cfg:       tt.fields.cfg,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			err := m.LoadAll()
			if (err != nil) != tt.wantErr {
				assert.Fail(fmt.Sprintf("Manager.LoadAll() error = %v, wantErr %v", err, tt.wantErr))
			}
		})
	}
}

func TestManager_List(t *testing.T) {
	type fields struct {
		log       *logrus.Logger
		cfg       *config.Config
		ghc       github.RepositoryService
		git       ackdevgit.OpenCloner
		repoCache []*Repository
	}
	type args struct {
		filters []Filter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Repository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				log:       tt.fields.log,
				cfg:       tt.fields.cfg,
				ghc:       tt.fields.ghc,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			if got := m.List(tt.args.filters...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_clone(t *testing.T) {
	type fields struct {
		log       *logrus.Logger
		cfg       *config.Config
		ghc       github.RepositoryService
		git       ackdevgit.OpenCloner
		repoCache []*Repository
	}
	type args struct {
		repoName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				log:       tt.fields.log,
				cfg:       tt.fields.cfg,
				ghc:       tt.fields.ghc,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			if err := m.clone(tt.args.repoName); (err != nil) != tt.wantErr {
				t.Errorf("Manager.clone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ensureFork(t *testing.T) {
	type fields struct {
		cfg       *config.Config
		ghc       github.RepositoryService
		git       ackdevgit.OpenCloner
		repoCache []*Repository
	}
	type args struct {
		repo *Repository
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				log:       tt.fields.log,
				cfg:       tt.fields.cfg,
				ghc:       tt.fields.ghc,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			if err := m.ensureFork(tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ensureFork() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ensureClone(t *testing.T) {
	type fields struct {
		log       *logrus.Logger
		cfg       *config.Config
		ghc       github.RepositoryService
		git       ackdevgit.OpenCloner
		repoCache []*Repository
	}
	type args struct {
		repo *Repository
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				log:       tt.fields.log,
				cfg:       tt.fields.cfg,
				ghc:       tt.fields.ghc,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			if err := m.ensureClone(tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("Manager.ensureClone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_EnsureAll(t *testing.T) {
	type fields struct {
		log       *logrus.Logger
		cfg       *config.Config
		ghc       github.RepositoryService
		git       ackdevgit.OpenCloner
		repoCache []*Repository
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				log:       tt.fields.log,
				cfg:       tt.fields.cfg,
				ghc:       tt.fields.ghc,
				git:       tt.fields.git,
				repoCache: tt.fields.repoCache,
			}
			if err := m.EnsureAll(); (err != nil) != tt.wantErr {
				t.Errorf("Manager.EnsureAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
