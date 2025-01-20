package main

import (
	"context"
)

func (g *GoLinkServer) GenerateShortURL(inputURL, hostname string) (string, error) {

	key, err := g.store.EncodeURL(inputURL)
	if err != nil {
		return "", err
	}
	// store the key value
	g.store.mapURL(context.Background(), key, inputURL)

	// the short form of the link
	return hostname + "/" + key, nil
}
