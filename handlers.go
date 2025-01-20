package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (g *GoLinkServer) shortenURL(w http.ResponseWriter, r *http.Request) {
	link := r.PostFormValue("fullUrl")
	returnHTMLFormat := `<form
      class="urlForm"
      hx-post="/make-it-smol"
      hx-target="form"
      hx-swap="outerHTML"
    >
      <input
        type="text"
        autocomplete="off"
        placeholder="Enter your URL to shorten!"
        name="fullUrl"
        id="linkInput"
				value="%s"
      />
      <button 
				type="button" 
				hx-on:click="navigator.clipboard.writeText(document.getElementById('linkInput').value);
				this.textContent ='Copied';"
				class="copy-btn">
				Copy
			</button>
			<button class="reset-btn" hx-get="/" hx-target="body">Reset</button>
    </form>`

	var inputHTML string

	// check if it's a valid link
	err := g.checkLink(link)
	if err != nil {
		http.Redirect(w, r, "/not-found", 301)
		return
	}

	shortUrl, err := g.GenerateShortURL(link, r.Host)
	if err != nil {
		log.Println(err)
		return
	}
	inputHTML = fmt.Sprintf(returnHTMLFormat, shortUrl)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(inputHTML))
}

func (g *GoLinkServer) checkLink(link string) error {
	client := &http.Client{
		Timeout: time.Second * 7,
	}
	_, err := client.Get(link)
	if err != nil {
		return err
	}
	return nil
}

func (g *GoLinkServer) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	urlKey := r.PathValue("id")
	mappedUrl, _ := g.store.getFullURL(r.Context(), urlKey)
	http.Redirect(w, r, mappedUrl, 302)
}
