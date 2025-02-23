/*
Package wiki is the page component for playing the game
*/
package wiki

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/bruceesmith/go-wikiracing/backend/api"
	"github.com/bruceesmith/go-wikiracing/frontend/actions"
	"github.com/bruceesmith/go-wikiracing/frontend/observables"
	"github.com/bruceesmith/logger"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

// Wiki represents the Wikipedia section of the game's web page
type Wiki struct {
	app.Compo
	ctx                  app.Context
	start, goal, current string
	Page                 string
	State                state
}

var (
	Default = Wiki{
		State: ready,
	}
	pageRe = regexp.MustCompile(`(?ms).+<body .+?>(.+)</body>`)
)

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

// Render is called by the UI goroutine each tine Wiki section of the web browser window
// needs an update
func (w *Wiki) Render() app.UI {
	return app.Div().Body(
		&theTopbar,
		&tmr,
		app.If(
			w.State == finished,
			func() app.UI {
				return app.Div().Body(
					app.Text("Goal reached!"),
				).
					Style("place-content", "center")
			},
		),
		app.Raw(w.Page),
	).
		Style("display", "grid").
		Style("place-content", "center")
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

// get fetches an HTML page from Wikipedia
func (w *Wiki) get(subject string) (s string, err error) {
	req := api.WikiPageRequest{Subject: subject}
	bites, err := json.Marshal(req)
	resp, err := http.Post("/api/wikipage", "application/json", bytes.NewBuffer(bites))
	if err != nil {
		logger.Error("Wiki.OnMount error fetching "+subject, "error", err.Error())
		return "", fmt.Errorf("Wiki.OnMount error fetching %s: [%w]", subject, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Wiki.OnMount error reading WikiPage response", "error", err.Error())
		return "", fmt.Errorf("Wiki.OnMount error reading WikiPage response: [%w]", err)
	}
	pageResponse := api.WikiPageResponse{}
	err = json.Unmarshal(body, &pageResponse)
	if err != nil {
		logger.Error("Wiki.OnMount error unmarshaing WikiPage response", "error", err.Error())
		return "", fmt.Errorf("Wiki.OnMount error unmarshaing WikiPage response: [%w]", err)
	}
	return pageResponse.Page, nil
}

// goalReached is an Action handler invoked when the "goal"
// Action is triggered
func (w *Wiki) goalReached(ctx app.Context) {
	ctx.SetState(observables.WikiState(), finished)
	tmr.finished()
}

// OnMount is called once, when the Wiki HTML is first added to the DOM
func (w *Wiki) OnMount(ctx app.Context) {
	// Register the action "pageloaded" which is triggered each
	// time a fresh Wkipedia page is fetched
	ctx.Handle(actions.PageLoaded(), w.updatePage)
	ctx.ObserveState(observables.WikiState(), &w.State)

	// Load the starting page in the background
	ctx.Async(
		func() {
			page, err := w.get("/wiki/" + w.start)
			if err != nil {
				return
			}
			w.current = w.start
			ctx.NewActionWithValue(actions.PageLoaded(), page)
		},
	)
	w.ctx = ctx
	// Register a Go function that is called from the Javascript
	// onclick event handler to analyze the href and potentially
	// fetch & load a fresh Wikipedia page
	app.Window().Set("wikiUrlClicked", app.FuncOf(w.urlclick))
}

// Targets sets the start and goal Wikipedia subjects
func (w *Wiki) Targets(start, goal string) {
	w.start = start
	w.goal = goal
}

// updatePage is an Action handler invoked when the "pageloaded"
// Action is triggered. It extracts the body of the Wikipedia HTML
// page and inserts an "onclick" attribute at each Wikipedia link
// in that HTML
func (w *Wiki) updatePage(ctx app.Context, a app.Action) {
	// Get the full page HTML source passed via the Action
	page, ok := a.Value.(string)
	if !ok {
		logger.Error("Wiki.updatePage internal error, unexpected type in Action.Value")
		return
	}
	// Extract the HTML between (but not including) the <body> and </body> tags
	matches := pageRe.FindSubmatch([]byte(page))
	if len(matches) == 0 {
		logger.Error("Wiki.updatePage got zero matches")
		return
	}
	// Update every <a> to include an onclick event trap
	body := strings.Replace(string(matches[1]), `<a `, `<a onclick="wikiAnchorClick(event)" `, -1)
	body = strings.Replace(body, `action="/w/index.php"`, `action="/wiki/Special:Search" onclick="wikiAnchorClick(event)"`, -1)
	// Update the Wiki Racing content
	w.Page = "<div>" + body + "</div>"
	if strings.ToLower(w.current) == "/wiki/"+strings.ToLower(w.goal) {
		w.goalReached(ctx)
	}
}

// urlclick is a method that is called from the Javascript
// onclick event handler to analyze the href and potentially
// fetch & load a fresh Wikipedia page
func (w *Wiki) urlclick(this app.Value, args []app.Value) (x any) {
	if len(args) == 0 {
		logger.Error("onclick handler received no arguments")
		return
	}
	href := args[0].String()
	url, err := url.Parse(href)
	if err != nil {
		logger.Error("cannot parse href", "href", href, "error", err.Error())
		return
	}
	if (strings.HasPrefix(url.Path, "/wiki/") || strings.HasPrefix(url.Path, "/static/")) && !strings.Contains(url.Path, "Special:Search") {
		// Load the requested page in the background
		w.ctx.Async(
			func() {
				page, err := w.get(url.Path)
				if err != nil {
					return
				}
				w.current = url.Path
				w.ctx.NewActionWithValue(actions.PageLoaded(), page)
			},
		)
	}
	return
}
