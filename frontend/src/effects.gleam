// Module effects contains all the Lustre Effects that are
// dispatched in the Update component of the Model-View-Update
// loop

import gleam/dynamic/decode
import gleam/json

import lustre/effect.{type Effect}
import rsvp

import endpoints.{type EP, type Goal, type Start, EP}
import msg.{type Msg}

/// dark_mode determines if OS/browser Dark mode is in effect
///
pub fn dark_mode_on(dark: fn(Bool) -> msg) -> Effect(msg) {
  effect.from(fn(dispatch) {
    let set = dark_mode()
    dispatch(dark(set))
  })
}

/// Interface to a JavaScript function which accesses the
/// window.matchMedia(() function
///
@external(javascript, "./app.ffi.mjs", "dark_mode")
fn dark_mode() -> Bool {
  False
}

const api_wiki_page = "/api/wikipage"

/// get_wiki_page fetches a Wiki page from the API server
///
/// Parameters:
///   subject: the Wiki topic to be fetched. This should
///            include the "/wiki/" path prefix
///   on_fetch: the Gleam Msg function that is called when the
///             Wiki page has been fetched
///
pub fn get_wiki_page(
  subject: String,
  on_fetch: fn(Result(String, rsvp.Error)) -> Msg,
) -> Effect(Msg) {
  let body = json.object([#("subject", json.string(subject))])
  let handler = rsvp.expect_text(on_fetch)

  rsvp.post(api_wiki_page, body, handler)
}

/// scroll_up positions a Wiki page to the top
///
pub fn scroll_up(on_scroll: msg) -> Effect(msg) {
  effect.from(fn(dispatch) {
    scroll_to_top()
    dispatch(on_scroll)
  })
}

/// Interface to the JavaScript window.scrollTo() function
///
@external(javascript, "./app.ffi.mjs", "scroll_to_top")
fn scroll_to_top() -> Nil {
  Nil
}

const api_special_random = "/api/specialrandom"

/// special_random calls the "specialrandom" API to
/// obtain two random Wiki topics
///
/// Parameters:
///   on_fetch: the Gleam Msg function that is called when the two
///             topics have been fetched
///
pub fn special_random(
  on_fetch: fn(Result(#(EP(Goal), EP(Start)), rsvp.Error)) -> Msg,
) -> Effect(Msg) {
  let decoder = {
    use random_goal <- decode.field("goal", decode.string)
    use random_start <- decode.field("start", decode.string)

    let goal: EP(Goal) = EP(random_goal)
    let start: EP(Start) = EP(random_start)

    decode.success(#(goal, start))
  }
  let handler = rsvp.expect_json(decoder, on_fetch)

  rsvp.get(api_special_random, handler)
}

/// start_timer fires a JavaScript timer and returns its ID
///
/// Parameters:
///   interval: the number of seconds between timer ticks
///   on_start: the Gleam Msg which is fired after the timer starts
///   on_tick: the Gleam Msg fired at each timer tick
///
pub fn start_timer(
  interval: Int,
  on_start: fn(Int) -> msg,
  on_tick: msg,
) -> Effect(msg) {
  effect.from(fn(dispatch) {
    let id = set_interval(interval, fn() { dispatch(on_tick) })

    dispatch(on_start(id))
  })
}

/// Interface to the JavaScript window.setInterval() function
///
@external(javascript, "./app.ffi.mjs", "set_interval")
fn set_interval(interval: Int, _cb: fn() -> a) -> Int {
  interval
}

/// stop_timer shuts down a JavaScript timer
///
/// Parameters:
///   id: a saved timer ID
///   on_stop: the Gleam Msg fired when the timer has stopped
///
pub fn stop_timer(id: Int, on_stop: msg) -> Effect(msg) {
  effect.from(fn(dispatch) {
    clear_interval(id)
    dispatch(on_stop)
  })
}

/// Interface to the JavaScript window.clearInterval() function
///
@external(javascript, "./app.ffi.mjs", "clear_interval")
fn clear_interval(_: Int) -> Nil {
  Nil
}
