import gleam/http
import gleam/http/request
import gleam/http/response
import gleam/httpc
import gleam/int
import gleam/result
import gleam/uri

const special_random_url = "https://en.wikipedia.org/wiki/Special:Random"

// 5 seconds timeout
const httpc_timeout = 5000

const wiki_url = "https://en.wikipedia.org"

/// Client is a record of functions that abstract the operations performed by the
/// client module.
///
pub type Client {
  Client(
    get_random: fn() -> Result(String, ClientError),
    get: fn(String) -> Result(WikiFileResult, ClientError),
  )
}

/// live returns a Client that uses the real Wikipedia API.
///
pub fn live() -> Client {
  Client(get_random: get_random, get: get)
}

/// ClientError is the set of possible errors returned by get_random() and get()
/// 
pub type ClientError {
  BadWikiResponse(String)
  HttpError(String)
  NoContentTypeHeader
  NoLocationHeader
  RequestCreationError
}

/// client_error_to_string converts a ClientError to a String
///
pub fn client_error_to_string(error: ClientError) -> String {
  case error {
    BadWikiResponse(message) -> "Bad Wikipedia response: " <> message
    HttpError(message) -> "HTTP error: " <> message
    NoContentTypeHeader -> "No content-type header found"
    NoLocationHeader -> "No location header found"
    RequestCreationError -> "Failed to create request"
  }
}

/// get_random returns the value of the Location header from a request to
/// https://en.wikipedia.org/wiki/Special:Random
///
pub fn get_random() -> Result(String, ClientError) {
  // Create a Request from the Special:Random URL
  use req <- result.try(
    request.to(special_random_url)
    |> result.map_error(fn(_) { RequestCreationError }),
  )

  // Make the request acceptable to Wikipedia
  let req =
    request.set_method(req, http.Head)
    |> request.set_header("Accept", "*/*")
    |> request.set_header("User-Agent", "wrspa/1.0")

  // Send the request and get the response
  let config =
    httpc.configure()
    |> httpc.timeout(httpc_timeout)
  use resp <- result.try(
    httpc.dispatch(config, req)
    |> result.map_error(fn(e) {
      HttpError("Failed to send request: " <> http_error_to_string(e))
    }),
  )

  // Try to extract the value of the Location header
  case resp.status {
    302 -> {
      use loc <- result.try(
        response.get_header(resp, "location")
        |> result.map_error(fn(_) { NoLocationHeader }),
      )
      use uri <- result.try(
        uri.parse(loc)
        |> result.map_error(fn(_) { BadWikiResponse(loc) }),
      )
      case uri.path {
        "/wiki/" <> rest -> Ok(rest)
        _ -> Error(BadWikiResponse(loc))
      }
    }
    _ -> {
      Error(BadWikiResponse(
        "Unexpected special:random status code: " <> int.to_string(resp.status),
      ))
    }
  }
}

/// WikiFileResult is the success value returned by get()
/// 
pub type WikiFileResult {
  WikiFileResult(body: String, content_type: String)
}

/// get fetches a file from Wikipedia. This includes both static files with
/// URL path prefixes of "/static/" and "/w/" and Wikipedia subject files
/// with URL path prefix of "/wiki"
/// 
pub fn get(path: String) -> Result(WikiFileResult, ClientError) {
  // Create a Request
  use req <- result.try(
    request.to(wiki_url <> path)
    |> result.map_error(fn(_) { RequestCreationError }),
  )

  // Make the request acceptable to Wikipedia
  let req =
    request.set_header(req, "Accept", "*/*")
    |> request.set_header("User-Agent", "wrspa/1.0")

  let config =
    httpc.configure()
    |> httpc.follow_redirects(True)
    |> httpc.timeout(httpc_timeout)

  // Send the request and get the response
  use resp <- result.try(
    httpc.dispatch(config, req)
    |> result.map_error(fn(e) {
      HttpError(
        "Failed to send request to "
        <> req.path
        <> ": "
        <> http_error_to_string(e),
      )
    }),
  )

  // Extract the body and the content-type
  case resp.status {
    200 -> {
      use content_type <- result.try(
        response.get_header(resp, "content-type")
        |> result.map_error(fn(_) { NoContentTypeHeader }),
      )
      Ok(WikiFileResult(resp.body, content_type))
    }
    _ -> {
      Error(BadWikiResponse(
        "Unexpected status code: " <> int.to_string(resp.status),
      ))
    }
  }
}

/// http_error_to_string converts an httpc.HttpError to a String
///
fn http_error_to_string(error: httpc.HttpError) -> String {
  case error {
    httpc.InvalidUtf8Response -> "Invalid UTF-8 response"
    httpc.FailedToConnect(ip4, ip6) ->
      "Failed to connect (IPv4: "
      <> connect_error_to_string(ip4)
      <> ", IPv6: "
      <> connect_error_to_string(ip6)
      <> ")"
    httpc.ResponseTimeout -> "Response timeout"
  }
}

/// connect_error_to_string converts an httpc.ConnectError to a String
///
fn connect_error_to_string(error: httpc.ConnectError) -> String {
  case error {
    httpc.Posix(code) -> code
    httpc.TlsAlert(code, detail) -> "TLS alert " <> code <> ": " <> detail
  }
}
