package main

import (
	"context"
	"encoding/base64"
	"fmt" //nolint:goimports
	"github.com/google/go-github/v47/github"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func main() {
	owner := "jameswoolfenden"
	repository := "terraform-aws-activemq"
	keyName := "HELLO"
	keyText := "hello world"
	response, err := SetRepoSecret(owner, repository, keyText, keyName)
	if err != nil {
		log.Print(response)
		log.Print(err)
	}
}

func SetRepoSecret(owner string, repository string, keyText string, keyName string) (*github.Response, error) {
	keyID, publicKey, err := getPublicKeyDetails(owner, repository)

	if err != nil {
		return nil, err
	}

	encryptedBytes, err := encryptPlaintext(keyText, publicKey)

	if err != nil {
		return nil, err
	}

	encryptedValue := base64.StdEncoding.EncodeToString(encryptedBytes)

	// Create an EncryptedSecret and encrypt the plaintext value into it
	eSecret := &github.EncryptedSecret{
		Name:           keyName,
		KeyID:          keyID,
		EncryptedValue: encryptedValue,
	}

	ctx, client := getGithubClient()

	response, err := client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repository, eSecret)

	if err != nil {
		return response, err
	}
	return response, nil
}

func getGithubClient() (context.Context, *github.Client) {
	token := os.Getenv("GITHUB_TOKEN")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client
}

func getPublicKeyDetails(owner, repository string) (keyID, pkValue string, err error) {
	ctx, client := getGithubClient()

	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repository)
	if err != nil {
		return keyID, pkValue, err
	}

	return publicKey.GetKeyID(), publicKey.GetKey(), err
}

func encryptPlaintext(plaintext, publicKeyB64 string) ([]byte, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyB64)
	if err != nil {
		return nil, err
	}

	var publicKeyBytes32 [32]byte
	copiedLen := copy(publicKeyBytes32[:], publicKeyBytes)
	if copiedLen == 0 {
		return nil, fmt.Errorf("could not convert publicKey to bytes")
	}

	plaintextBytes := []byte(plaintext)
	var encryptedBytes []byte

	cipherText, err := box.SealAnonymous(encryptedBytes, plaintextBytes, &publicKeyBytes32, nil)
	if err != nil {
		return nil, err
	}

	return cipherText, nil
}
