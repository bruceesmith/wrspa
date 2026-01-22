import client
import gleam/http/response
import gleeunit
import gleeunit/should

pub fn main() -> Nil {
  gleeunit.main()
}

pub fn client_error_to_string_test() {
  client.client_error_to_string(client.BadWikiResponse("500"))
  |> should.equal("Bad Wikipedia response: 500")

  client.client_error_to_string(client.HttpError("Timeout"))
  |> should.equal("HTTP error: Timeout")

  client.client_error_to_string(client.NoContentTypeHeader)
  |> should.equal("No content-type header found")

  client.client_error_to_string(client.NoLocationHeader)
  |> should.equal("No location header found")

  client.client_error_to_string(client.RequestCreationError)
  |> should.equal("Failed to create request")
}

pub fn get_random_test() {
  // Mock a successful 302 redirect from Special:Random
  let mock_dispatch = fn(_config, _req) {
    Ok(
      response.new(302)
      |> response.set_header(
        "location",
        "https://en.wikipedia.org/wiki/Random_Page",
      ),
    )
  }

  let result = client.get_random(mock_dispatch)
  should.be_ok(result)
  let assert Ok(res) = result
  should.equal(res, "Random_Page")
}

pub fn get_test() {
  // Mock a successful 200 OK response
  let mock_dispatch = fn(_config, _req) {
    Ok(
      response.new(200)
      |> response.set_header("content-type", "text/html")
      |> response.set_body("<html><body>Test Body</body></html>"),
    )
  }

  let result = client.get("/wiki/Main_Page", mock_dispatch)
  should.be_ok(result)
  let assert Ok(res) = result
  should.equal(res.body, "<html><body>Test Body</body></html>")
  should.equal(res.content_type, "text/html")
}

pub fn get_404_test() {
  // Mock a 404 Not Found response
  let mock_dispatch = fn(_config, _req) {
    Ok(response.new(404) |> response.set_body("Not Found"))
  }

  let result = client.get("/wiki/ThisPageDoesNot_Exist_123", mock_dispatch)
  should.be_error(result)

  case result {
    Error(client.BadWikiResponse("Unexpected status code: 404")) -> Nil
    _ -> should.fail()
  }
}

pub fn get_invalid_url_test() {
  // Mock dispatcher (should not be called)
  let mock_dispatch = fn(_config, _req) { Ok(response.new(200)) }

  // Test invalid URL creation
  // Passing a control character like newline should fail URL parsing
  let result = client.get("\n", mock_dispatch)
  should.be_error(result)
  should.equal(result, Error(client.RequestCreationError))
}
