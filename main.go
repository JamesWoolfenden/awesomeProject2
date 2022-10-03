package main

import (
	"context"
	"github.com/google/go-github/v47/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	james := github.RepositoryListOptions{"Public", "owner", "", "", "", github.ListOptions{Page: 4, PerPage: 50}}
	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, "", &james)
	if err != nil {
		log.Fatal(err)
	}

	for _, repo := range repos {
		test := *repo
		log.Printf(test.String())
	}

	//temp := github.ActionsService{client}
	temp2, _, _ := client.Actions.ListRepoSecrets(ctx, "jameswoolfenden", "terraform-aws-activemq", nil)
	secrets := temp2.Secrets
	for _, secret := range secrets {
		temp3 := *secret
		log.Printf(temp3.Name)
	}

	client.Actions.CreateOrUpdateRepoSecret(ctx, "jameswoolfenden", "terraform-aws-activemq", "")
}
