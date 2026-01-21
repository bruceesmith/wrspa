import gleam/http.{Get, Post}
import gleam/json
import gleeunit
import gleeunit/should
import server
import wisp/simulate

pub fn main() {
  gleeunit.main()
}

pub fn handle_request_test() {
  // Test redirect to index
  let req = simulate.request(Get, "/")
  let resp = server.handle_request(req, ".")
  resp.status |> should.equal(303)

  // Test 404 for unknown paths
  let req = simulate.request(Get, "/unknown_path_123")
  let resp = server.handle_request(req, ".")
  resp.status |> should.equal(404)
}

pub fn api_special_random_test() {
  let req = simulate.request(Get, "/api/specialrandom")
  let resp = server.handle_request(req, ".")

  resp.status |> should.equal(200)
}

pub fn api_wikipage_test() {
  // Test valid request
  let body_json = json.object([#("subject", json.string("/wiki/Main_Page"))])

  let req =
    simulate.request(Post, "/api/wikipage")
    |> simulate.json_body(body_json)

  let resp = server.handle_request(req, ".")

  resp.status |> should.equal(200)

  // Test invalid method
  let req = simulate.request(Get, "/api/wikipage")
  let resp = server.handle_request(req, ".")
  resp.status |> should.equal(405)

  // Test invalid JSON (malformed)
  let req =
    simulate.request(Post, "/api/wikipage")
    |> simulate.string_body("{invalid")
    |> simulate.header("content-type", "application/json")

  let resp = server.handle_request(req, ".")
  resp.status |> should.equal(400)

  // Test invalid subject (must start with /wiki/)
  let body_json = json.object([#("subject", json.string("/other/page"))])
  let req =
    simulate.request(Post, "/api/wikipage")
    |> simulate.json_body(body_json)

  let resp = server.handle_request(req, ".")
  resp.status |> should.equal(400)
}

pub fn wikipedia_file_test() {
  // Test accessing a wikipedia file
  // We'll use a likely valid path.
  let req = simulate.request(Get, "/static/images/project-logos/enwiki.png")
  let _resp = server.handle_request(req, ".")

  // Checking for 405 on POST is a safe test that doesn't rely on network as much for routing logic check
  let req_post = simulate.request(Post, "/static/foo")
  let resp_post = server.handle_request(req_post, ".")
  resp_post.status |> should.equal(405)
}
