import argv
import glint

import daemon

pub fn main() -> Nil {
  glint.new()
  |> glint.with_name("wrserver")
  |> glint.global_help("Wiki Racing Server")
  |> glint.pretty_help(glint.default_pretty_help())
  |> glint.add(at: ["start"], do: daemon.daemon())
  |> glint.run(argv.load().arguments)
}