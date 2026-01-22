# Wiki Racing Server (Backend)

This is the backend service for the Wiki Racing application, written in Gleam. It serves as both a static file server for the frontend SPA and a proxy API to interact with Wikipedia, handling tasks like fetching random pages and retrieving page content without CORS issues.

## Architecture

The application is structured as follows:

*   **Entry Point (`src/wrserver.gleam`):** Uses `glint` to handle command-line arguments and dispatch the `start` command.
*   **Daemon (`src/daemon.gleam`):** Manages the application lifecycle using an OTP supervisor. It validates configuration (port, static directory) and keeps the server running.
*   **Web Server (`src/server.gleam`):** Built with `mist` and `wisp`.
    *   Serves static frontend assets from the configured directory.
    *   **`/api/specialrandom`:** Fetches two random Wikipedia page titles (start and goal).
    *   **`/api/wikipage`:** Proxies a Wikipedia page request, extracting the `<body>` content to return a cleaner HTML fragment.
    *   **`/static/` & `/w/`:** Proxies static assets (images, styles) directly from Wikipedia.
*   **Client (`src/client.gleam`):** A `gleam_httpc` wrapper for communicating with `en.wikipedia.org`. Handles HTTP requests, redirects, and error mapping.

## Building and Running

### Prerequisites
*   Gleam
*   Erlang/OTP

### Commands

**Build:**
```bash
gleam build
```

**Run (Development):**
To start the server, you need to use the `start` subcommand. You can optionally specify the port and the directory for static assets.
```bash
# Default (Port 8080, static dir "dist")
gleam run -- start

# Custom configuration
gleam run -- start --port 3000 --static /path/to/frontend/dist
```

**Testing:**
```bash
gleam test
```

**Formatting:**
```bash
gleam format
```

## Configuration

*   **Port:** Defaults to `8080`. configured via `--port`.
*   **Static Directory:** Defaults to `dist`. Configured via `--static`.
*   **Target:** Runs on the Erlang VM (`target = "erlang"` in `gleam.toml`).

## Key Dependencies

*   **`mist`**: HTTP server.
*   **`wisp`**: Web framework for request handling.
*   **`glint`**: CLI argument parsing.
*   **`gleam_httpc`**: HTTP client for fetching Wikipedia data.
*   **`gleam_otp`**: Supervisor and actor primitives.
*   **`gleam_json`**: JSON encoding/decoding.
