import gleam/erlang/process
import gleam/int
import gleam/io
import gleam/otp/actor
import gleam/otp/static_supervisor as supervisor
import gleam/otp/supervision
import gleam/result
import glint
import snag

import server
import simplifile

pub fn daemon() -> glint.Command(Nil) {
  use <- glint.command_help("Starts the Wiki Racing server")

  use port_flag <- glint.flag(
    glint.int_flag("port")
    |> glint.flag_default(8080)
    |> glint.flag_constraint(validate_port),
  )
  use static_flag <- glint.flag(
    glint.string_flag("static")
    |> glint.flag_default("dist")
    |> glint.flag_constraint(validate_static),
  )

  use _named_args, _args, flags <- glint.command()
  let assert Ok(port) = port_flag(flags)
  let assert Ok(static) = static_flag(flags)

  let assert Ok(_pid) = start_server_supervisor(port, static)

  io.println("✅ Daemon is running. Press Ctrl+C twice to stop.")

  process.sleep_forever()
}

pub fn start_server_supervisor(
  port: Int,
  static: String,
) -> Result(process.Pid, actor.StartError) {
  io.println("🚀 Starting daemon with static supervisor...")

  let worker_child = supervision.worker(fn() { server.serve(port, static) })

  supervisor.new(supervisor.OneForOne)
  |> supervisor.add(worker_child)
  |> supervisor.start
  |> result.map(fn(started) { started.pid })
}

pub fn validate_port(i: Int) -> Result(Int, snag.Snag) {
  case i <= 1024 || i > 65_535 {
    True ->
      snag.error(
        "value " <> int.to_string(i) <> " must be between 1025 and 65535",
      )
    False -> Ok(i)
  }
}

pub fn validate_static(s) -> Result(String, snag.Snag) {
  case simplifile.is_directory(s) {
    Ok(True) ->
      case simplifile.read_directory(s) {
        Ok(_) -> Ok(s)
        Error(_) -> snag.error("directory " <> s <> " is not readable")
      }
    Ok(False) -> snag.error(s <> " does not exist or is not a directory")
    Error(_) -> snag.error("unable to access " <> s)
  }
}
