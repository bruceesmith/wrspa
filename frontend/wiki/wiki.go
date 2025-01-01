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

	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/go-wikiracing/backend/api"
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
}

var (
	Default Wiki
	pageRe  = regexp.MustCompile(`(?ms).+<body .+?>(.+)</body>`)
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
		app.Div().Text("I am the wiki "+w.start+" "+w.goal),
		// app.Div().Body(
		// 	app.Raw(`<a href="https://www.w3schools.com" onclick="wikiAnchorClick(event)">Visit W3Schools.com!</a>`),
		// ),
		app.Raw(w.Page),
	)
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

// get fetches an HTML page from Wikipedia
func (w *Wiki) get(subject string) (s string) {
	var err error
	req := api.WikiPageRequest{Subject: subject}
	bites, err := json.Marshal(req)
	resp, err := http.Post("/api/wikipage", "application/json", bytes.NewBuffer(bites))
	if err != nil {
		logger.Error("Wiki.OnMount error fetching "+subject, "error", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Wiki.OnMount error reading WikiPage response", "error", err.Error())
		return
	}
	pageResponse := api.WikiPageResponse{}
	err = json.Unmarshal(body, &pageResponse)
	if err != nil {
		logger.Error("Wiki.OnMount error unmarshaing WikiPage response", "error", err.Error())
		return
	}
	return pageResponse.Page
}

// goalReached is an Action handler invoked when the "goal"
// Action is triggered
func (w *Wiki) goalReached(_ app.Context, _ app.Action) {
	fmt.Println("goal, goal, goal !!!")
}

// OnMount is called once, when the Wiki HTML is first added to the DOM
func (w *Wiki) OnMount(ctx app.Context) {
	// Register the action "pageloaded" which is triggered each
	// time a fresh Wkipedia page is fetched
	ctx.Handle("pageloaded", w.updatePage)
	// Register the action "goal" which is triggered if there
	// is a click on a Wiki link to the goal subject
	ctx.Handle("goal", w.goalReached)
	// Load the starting page in the background
	ctx.Async(
		func() {
			page := w.get(w.start)
			w.current = w.start
			ctx.NewActionWithValue("pageloaded", page)
		},
	)
	w.ctx = ctx
	// Register a Go function that is called from the Javascript
	// onclick event handler to analyze the href and potentially
	// fetch & load a fresh Wikipedia page
	app.Window().Set("wikiUrlClicked", app.FuncOf(w.urlclick))
}

// Targets sets the start and goal Wikipedia subjects
func (t *Wiki) Targets(start, goal string) {
	t.start = start
	t.goal = goal
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
	// Update every <a> to include an onclick event trap
	body := strings.Replace(string(matches[1]), `<a `, `<a onclick="wikiAnchorClick(event)" `, -1)
	body = strings.Replace(body, `action="/w/index.php"`, `action="/wiki/Special:Search" onclick="wikiAnchorClick(event)"`, -1)
	// Update the Wiki Racing content
	w.Page = "<div>" + body + "</div>"
	if strings.ToLower(w.current) == strings.ToLower(w.goal) {
		fmt.Println("goal! goal! goal!")
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
	fmt.Println("urlclick got href", href)
	url, err := url.Parse(href)
	if err != nil {
		logger.Error("cannot parse href", "href", href, "error", err.Error())
		return
	}
	if strings.HasPrefix(url.Path, "/wiki/") && !strings.Contains(url.Path, "Special:Search") {
		fmt.Println("I'll navigate to", href)
		// Load the requested page in the background
		w.ctx.Async(
			func() {
				subject := strings.TrimPrefix(url.Path, "/wiki/")
				page := w.get(subject)
				w.current = subject
				w.ctx.NewActionWithValue("pageloaded", page)
			},
		)
	} else {
		fmt.Println("na na, you cannot go to", href)
	}
	return
}
