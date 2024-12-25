package main

import (
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

func main() {
	app.Route("/", func() app.Composer { return newGame() })
	app.RunWhenOnBrowser()
	http.Handle("/", &app.Handler{
		Name:        "WikiRacing",
		Description: "A wiki racing game",
		Styles: []string{
			"/web/game.css",
		},
		RawHeaders: []string{
			`<script>
				function wikiAnchorClick(event, obj) {
  					event.preventDefault();
					wikiUrlClicked(obj);
				}
			</script>`,
		},
	})
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
