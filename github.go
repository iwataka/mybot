package main

import (
	"net/http"

	"github.com/google/go-github/github"
)

// GitHubAPI is a wrapper of github.Client
type GitHubAPI struct {
	*github.Client
	cache *MybotCache
}

// NewGitHubAPI creates an instance of GitHubAPI
func NewGitHubAPI(c *http.Client, cache *MybotCache) *GitHubAPI {
	return &GitHubAPI{
		github.NewClient(c),
		cache,
	}
}

// GitHubProject represents a GitHub project
type GitHubProject struct {
	User string `toml:"user"`
	Repo string `toml:"repo"`
}

func (a *GitHubAPI) LatestCommit(p GitHubProject) (*github.RepositoryCommit, error) {
	commits, _, err := a.Repositories.ListCommits(p.User, p.Repo, nil)
	if err != nil {
		return nil, err
	}
	latest := commits[0]
	userMap, userExists := a.cache.LatestCommitSHA[p.User]
	if !userExists {
		cache.LatestCommitSHA[p.User] = make(map[string]string)
	}
	sha, exists := userMap[p.Repo]
	if !exists || sha != *latest.SHA {
		cache.LatestCommitSHA[p.User][p.Repo] = *latest.SHA
		return latest, nil
	}
	return nil, nil
}
