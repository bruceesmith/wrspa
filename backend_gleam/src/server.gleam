import gleam/dynamic/decode
import gleam/http.{Get, Post}
import gleam/json
import gleam/list
import gleam/result
import gleam/string
import mist
import wisp.{type Request, type Response}
import wisp/wisp_mist

import client

/// serve starts the Mist web server
/// 
pub fn serve(port: Int, static_directory: String) {
  // This sets the logger to print INFO level logs, and other sensible defaults
  // for a web application.
  wisp.configure_logger()

  // Here we generate a secret key, but in a real application you would want to
  // load this from somewhere so that it is not regenerated on every restart.
  let secret_key_base = wisp.random_string(64)

  // Start the Mist web server.
  wisp_mist.handler(
    fn(req) { handle_request(req, static_directory) },
    secret_key_base,
  )
  |> mist.new
  |> mist.port(port)
  |> mist.start
}

/// handle_request handles incoming requests to the Mist web server
/// 
fn handle_request(req: Request, static_directory: String) -> Response {
  // Log information about the request and response.
  use <- wisp.log_request(req)

  // Return a default 500 response if the request handler crashes.
  use <- wisp.rescue_crashes

  // Rewrite HEAD requests to GET requests and return an empty body.
  use req <- wisp.handle_head(req)

  // Handle explicit file requests (e.g. /wrspa.js, /style.css)
  use <- wisp.serve_static(req, under: "/", from: static_directory)

  // Handle requests to the REST API and for static Wikipedia pages
  case wisp.path_segments(req) {
    [] -> wisp.redirect("/index.html")

    ["api", function, ..] -> api(req, function)

    ["static", ..] | ["w", ..] -> wikipedia_file(req)

    _ -> wisp.not_found()
  }
}

/// api provides the REST interface for the SPA
///
fn api(req: Request, function: String) -> Response {
  case req.method {
    Get if function == "specialrandom" -> special_random(req)

    Post if function == "wikipage" -> wikipage(req)

    _ -> wisp.method_not_allowed([Get, Post])
  }
}

/// special_random fetches the subjects for two random Wikipedia pages
///
pub fn special_random(req: Request) -> Response {
  use <- wisp.require_method(req, Get)

  // Fetch one random Wikipedia subject
  let result = {
    use first <- result.try(
      client.get_random()
      |> result.map_error(fn(e) {
        wisp.internal_server_error()
        |> wisp.html_body(
          "could not fetch any random endpoints: "
          <> client.client_error_to_string(e),
        )
      }),
    )

    // Fetch a second random Wikipedia subject
    use second <- result.try(
      client.get_random()
      |> result.map_error(fn(e) {
        wisp.internal_server_error()
        |> wisp.html_body(
          "could only fetch one random endpoint ("
          <> first
          <> "): "
          <> client.client_error_to_string(e),
        )
      }),
    )

    // Combine the two random subjects into a JSON response & return it
    let jason =
      json.object([
        #("start", json.string(first)),
        #("goal", json.string(second)),
      ])
    Ok(wisp.json_response(json.to_string(jason), 200))
  }

  case result {
    Ok(response) | Error(response) -> response
  }
}

/// wikipage fetches a Wikipedia page
/// 
fn wikipage(req: Request) -> Response {
  use <- wisp.require_method(req, Post)
  use jason <- wisp.require_json(req)

  let result = {
    // Decode the JSON body
    use WikiSubject(subject) <- result.try(
      decode.run(jason, wikipage_decoder())
      |> result.map_error(fn(errors) {
        wisp.bad_request(
          "Unable to decode request: " <> decode_errors_to_string(errors),
        )
      }),
    )

    // Make sure the subject starts with "/wiki/"
    use _ <- result.try(case string.starts_with(subject, "/wiki/") {
      True -> Ok(Nil)
      False ->
        Error(wisp.bad_request(
          "subject for api/wikipage must start with /wiki/",
        ))
    })

    // Fetch the page
    use file <- result.try(
      client.get(subject)
      |> result.map_error(fn(e) {
        wisp.internal_server_error()
        |> wisp.html_body(client.client_error_to_string(e))
      }),
    )

    // Extract the content between <body> and </body>
    use body <- result.try(
      body_content(file.body)
      |> result.map_error(fn(e) {
        wisp.internal_server_error()
        |> wisp.html_body(e)
      }),
    )

    Ok(
      wisp.ok()
      |> wisp.set_header("Content-Type", file.content_type)
      |> wisp.string_body(body),
    )
  }

  case result {
    Ok(response) | Error(response) -> response
  }
}

/// body_content extracts the content between <body> and </body>
/// 
pub fn body_content(page: String) -> Result(String, String) {
  // Get the content following the partial <body tag
  use #(_, back) <- result.try(
    string.split_once(page, "<body")
    |> result.replace_error("cannot locate <body"),
  )

  // Get the content preceding the </body> tag
  use #(body_front, _) <- result.try(
    string.split_once(back, "</body>")
    |> result.replace_error("cannot locate </body>"),
  )

  // Split off the remainder of the leading body tag (i.e. its attributes)
  use #(_, body) <- result.try(
    string.split_once(body_front, ">")
    |> result.replace_error("cannot locate </body>"),
  )
  Ok(body)
}

/// WikiSubject is the JSON request to the api/wikipage endpoint
/// 
type WikiSubject {
  WikiSubject(subject: String)
}

/// wikipage_decoder decodes a WikiSubject from JSON
///
fn wikipage_decoder() -> decode.Decoder(WikiSubject) {
  use subject <- decode.field("subject", decode.string)
  decode.success(WikiSubject(subject))
}

/// decode_errors_to_string converts a list of DecodeError to a String
///
fn decode_errors_to_string(errors: List(decode.DecodeError)) -> String {
  errors
  |> list.map(fn(error) {
    let path = string.join(error.path, ".")
    "Expected " <> error.expected <> ", found " <> error.found <> " at " <> path
  })
  |> string.join(", ")
}

/// wikipedia_file fetches static Wikipedia files
/// 
fn wikipedia_file(req: Request) -> Response {
  use <- wisp.require_method(req, Get)

  let result = {
    let segments = wisp.path_segments(req)

    // Make sure the path has at least two segments
    use _ <- result.try(case segments {
      [] | [_] ->
        Error(
          wisp.internal_server_error() |> wisp.html_body("no file specified"),
        )
      _ -> Ok(Nil)
    })

    // Prepare the full path & fetch the file
    let path = "/" <> string.join(segments, with: "/")
    use file <- result.try(
      client.get(path)
      |> result.map_error(fn(e) {
        wisp.internal_server_error()
        |> wisp.html_body(client.client_error_to_string(e))
      }),
    )

    // Send back the file
    Ok(
      wisp.ok()
      |> wisp.set_header("Content-Type", file.content_type)
      |> wisp.string_body(file.body),
    )
  }

  case result {
    Ok(response) | Error(response) -> response
  }
}
