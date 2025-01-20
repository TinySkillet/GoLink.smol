package main

func (g *GoLinkServer) GenerateShortURL(inputURL, hostname string) (string, error) {

	key, err := g.store.EncodeURL(inputURL)
	if err != nil {
		return "", err
	}
	// the short form of the link
	return hostname + "/" + key, nil
}
