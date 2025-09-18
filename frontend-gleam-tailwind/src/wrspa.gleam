//// Module wrspa contains the main Gleam program for the
//// Wiki Racing single-page application

import lustre

import init.{init}
import update.{update}
import view.{view}

// -----------------------------------------------------------------------------
//
// MAIN ------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

pub fn main() {
  let app = lustre.application(init, update, view)
  let assert Ok(_) = lustre.start(app, "#app", Nil)

  Nil
}
