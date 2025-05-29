# Wiki Racing SPA and API Server

## About

The Wiki Racing SPA is a simple single-player implementation of the 
[Wiki Game](https://en.wikipedia.org/wiki/Wikipedia:Wiki_Game) in which a player endeavours to navigate from
a starting Wkipedia page to a goal page by clicking only on wikilinks. The objective is to reach the goal in
the shortest time while clicking on the smallest numbert of wikilinks.

## Components

There are two components
- a [Single-Page Application](https://en.wikipedia.org/wiki/Single-page_application) executing in a 
  web browser
- a [REST API server](https://en.wikipedia.org/wiki/REST)

## SPA

Implemented in [Gleam](https://gleam.run/) using the [Lustre web framework](https://github.com/lustre-labs/lustre), the SPA is delivered as a single minified JavaScript bundle. A player LOADS the game by simply
visiting a URL which responds with a simple HTML index page.

## API server

The REST API server is a small [Go](https://go.dev/) binary that can run on any platform supported by Go 
(including Windows, Linux, macOS) or hosted on any Docker platform.

## Building from source

1. Clone the repository 
  - `git clone https://github.com/bruceesmith/wrspa.git`
  - `cd wrspa`
2. Install Go,  Gleam and Tailwind CSS
3. Build the SPA
  - `cd frontend`
  - `make build-prod`
  - `cd ..`
4. Build the server
  - `cd backend`
  - `make build`
  - `cd ..`

