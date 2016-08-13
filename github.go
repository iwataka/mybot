package main

import (
	"net/http"

	"github.com/google/go-github/github"
)

type GitHubAPI struct {
	*github.Client
	cache *MybotCache
}

func NewGitHubAPI(c *http.Client, cache *MybotCache) *GitHubAPI {
	return &GitHubAPI{
		github.NewClient(c),
		cache,
	}
}

type GitHubProject struct {
	User string
	Repo string
}

func (a *GitHubAPI) LatestCommit(p GitHubProject) (*github.RepositoryCommit, error) {
	commits, _, err := a.Repositories.ListCommits(p.User, p.Repo, nil)
	if err != nil {
		return nil, err
	}
	latest := commits[0]
	userMap, userExists := a.cache.LatestCommitSHA[p.User]
	sha := ""
	repoExists := false
	if userExists {
		sha, repoExists = userMap[p.Repo]
	}
	if !userExists || !repoExists || sha != *latest.SHA {
		msg := p.User + "/" + p.Repo + "\n" + *latest.HTMLURL
		_, err := twitterAPI.PostTweet(msg, nil)
		if err != nil {
			return nil, err
		}
		if !userExists {
			a.cache.LatestCommitSHA[p.User] = make(map[string]string)
		}
		a.cache.LatestCommitSHA[p.User][p.Repo] = *latest.SHA
		return latest, nil
	}
	return nil, nil
}
