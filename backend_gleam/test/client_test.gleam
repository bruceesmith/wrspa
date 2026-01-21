import client
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
  // This is an integration test that hits Wikipedia.
  // It may fail if there is no internet connection.
  let result = client.get_random()
  should.be_ok(result)
  let assert Ok(res) = result
  should.not_equal(res, "")
}

pub fn get_test() {
  // This is an integration test that hits Wikipedia.
  // We'll try to fetch the Main Page which should exist.
  let result = client.get("/wiki/Main_Page")
  should.be_ok(result)
  let assert Ok(res) = result
  should.not_equal(res.body, "")
  should.not_equal(res.content_type, "")
}

pub fn get_404_test() {
  // Test 404 error from Wikipedia
  let result = client.get("/wiki/ThisPageDoesNot_Exist_123")
  should.be_error(result)
  
  case result {
    Error(client.BadWikiResponse(_)) -> Nil
    _ -> should.fail()
  }
}

pub fn get_invalid_url_test() {
  // Test invalid URL creation
  // Passing a control character like newline should fail URL parsing
  let result = client.get("\n")
  should.be_error(result)
  should.equal(result, Error(client.RequestCreationError))
}
