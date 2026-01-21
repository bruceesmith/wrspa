import daemon
import gleam/erlang/process
import gleeunit
import gleeunit/should
import glint
import simplifile

pub fn main() -> Nil {
  gleeunit.main()
}

pub fn daemon_test() {
  // Test that daemon returns a command that fails with invalid flags
  glint.new()
  |> glint.add(at: ["start"], do: daemon.daemon())
  |> glint.execute(["start", "--port", "80"])
  |> should.be_error

  glint.new()
  |> glint.add(at: ["start"], do: daemon.daemon())
  |> glint.execute(["start", "--static", "non_existent_dir_12345"])
  |> should.be_error
}

pub fn start_server_supervisor_test() {
  // Test starting the supervisor with valid parameters
  let result = daemon.supervisor(8081, "src")
  result |> should.be_ok

  let assert Ok(pid) = result
  // Clean up by sending a normal exit signal to the supervisor process
  process.send_exit(to: pid)
}

pub fn validate_static_test() {
  // Test with non-existent directory
  daemon.validate_static("non_existent_dir_12345")
  |> should.be_error

  // Test with existing directory
  daemon.validate_static("src")
  |> should.be_ok

  // Test with a file (should fail)
  let file = "test_file_for_validate_static.txt"
  // Create a dummy file
  let assert Ok(_) = simplifile.create_file(file)

  // It should be an error because it is a file, not a directory
  daemon.validate_static(file)
  |> should.be_error

  // Clean up
  let assert Ok(_) = simplifile.delete(file)
}

pub fn validate_port_test() {
  // Test valid ports
  daemon.validate_port(1025)
  |> should.be_ok

  daemon.validate_port(8080)
  |> should.be_ok

  daemon.validate_port(65_535)
  |> should.be_ok

  // Test invalid ports (too low)
  daemon.validate_port(80)
  |> should.be_error

  daemon.validate_port(1024)
  |> should.be_error

  // Test invalid ports (too high)
  daemon.validate_port(65_536)
  |> should.be_error

  daemon.validate_port(100_000)
  |> should.be_error
}
