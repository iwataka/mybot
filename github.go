package main

import "github.com/google/go-github/github"

var githubClient = github.NewClient(nil)

func githubCommit(user, repo string) (*github.RepositoryCommit, error) {
	commits, _, err := githubClient.Repositories.ListCommits(user, repo, nil)
	if err != nil {
		return nil, err
	}
	latest := &commits[0]
	userMap, userExists := cache.LatestCommitSHA[user]
	sha := ""
	repoExists := false
	if userExists {
		sha, repoExists = userMap[repo]
	}
	if !userExists || !repoExists || sha != *latest.SHA {
		return latest, nil
		msg := user + "/" + repo + "\n" + *latest.HTMLURL
		_, err := twitterApi.PostTweet(msg, nil)
		if err != nil {
			return nil, err
		}
		if !userExists {
			cache.LatestCommitSHA[user] = make(map[string]string)
		}
		cache.LatestCommitSHA[user][repo] = *latest.SHA
	}
	return nil, nil
}
