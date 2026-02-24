//// Module view provides the view() function, one of the
//// Lustre Model-View-Update triumvirate, and a series of helper
//// functions

import gleam/bool.{or}
import gleam/dynamic/decode
import gleam/int
import gleam/option.{type Option, None, Some}
import gleam/string

import lustre/attribute.{attribute, class, disabled, for, id, required}
import lustre/element.{type Element}
import lustre/element/html as h
import lustre/event

import m3e/button
import m3e/form_field

import endpoints.{
  type EP, type Endpoints, type Goal, type Start, EP, active_goal, custom_goal,
  custom_start, random_goal, random_start,
}

import model.{
  type Model, type State, ChoosingGame, Completed, CustomGame, Paused, Playing,
  RandomGame, ReadyToPlay,
}
import msg.{
  type Msg, Click, CustomEndPointsSelected, CustomGoalChanged, CustomSelected,
  CustomStartChanged, GamePaused, GameResumed, GameStarted, NavigateBack,
  NavigateForward, NewGame, RandomEndPointsDisplayed, RandomSelected,
  RedrawRandom, RestartGame,
}
import navigation.{type Navigation, navigation_possible}

// -----------------------------------------------------------------------------
//
// VIEW ------------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

/// view is called in the View phase of Model-View-Update after update()
/// has processed a Nsg and possibly created a new Mode. view() generates
/// new HTMl based on the Model that is then displayed
///
/// Parameters:
///   model: the current SPA state, an instance of the Model record
///
pub fn view(model: Model) -> Element(Msg) {
  let body = case model.state {
    ChoosingGame | RandomGame | CustomGame ->
      choosing(
        model.state,
        model.endpoints,
        model.goal_error,
        model.start_error,
        model.rsvp_error,
      )
    ReadyToPlay | Playing | Paused | Completed -> playing(model)
  }
  let title =
    h.div([class("grid grid-rows-1 border-none")], [
      h.p(
        [
          class(
            "font-bold italic pb-1 pt-1 justify-self-center text-primary text-3xl",
          ),
        ],
        [h.text("Wiki Racing")],
      ),
    ])
  h.div([], [title, ..body])
}

// -----------------------------------------------------------------------------
//
// SETUP PAGE and HELPERS------------------------------------------------------------
//
// -----------------------------------------------------------------------------

/// choosing builds the setup phase page
///
/// Parameters:
///   state: current game state (one of the setup phase values)
///   endpoints: the games' start and goal
///   goal_error: error in the requested Goal
///   start_error: error in the request Start
///   rsvp_error: error interacting with or from the API server
///
fn choosing(
  state: State,
  endpoints: Endpoints,
  goal_error: Option(String),
  start_error: Option(String),
  rsvp_error: Option(String),
) -> List(Element(Msg)) {
  let second_row = case state {
    ChoosingGame -> element.none()
    RandomGame ->
      random(endpoints |> random_start, endpoints |> random_goal, rsvp_error)
    CustomGame ->
      custom(
        endpoints |> custom_start,
        endpoints |> custom_goal,
        goal_error,
        start_error,
      )
    _ -> element.none()
  }
  [
    h.div([class("grid grid-cols-2 gap-4 justify-center")], [
      button.render(button.new("Random", button.Filled), [
        class("justify-self-end"),
        event.on_click(RandomSelected),
      ]),
      button.render(button.new("Custom", button.Filled), [
        class("justify-self-start"),
        event.on_click(CustomSelected),
      ]),
    ]),
    second_row,
  ]
}

/// custom builds the section of the setup page which handles
/// user selection of Start and Goal
///
/// Parameters:
///   start: the selected starting Wiki topic
///   goal: the selected Wiki goal topic
///   goal_error: error in the requested Goal
///   start_error: error in the request Start
///
fn custom(
  start: EP(Start),
  goal: EP(Goal),
  goal_error: Option(String),
  start_error: Option(String),
) -> Element(Msg) {
  let #(error_line, _) = custom_error_line(start_error, goal_error)
  let EP(gl): EP(Goal) = goal
  let EP(st): EP(Start) = start

  h.div([class("grid grid-rows-[1fr_2fr_1fr] gap-2")], [
    h.div([class("self-center justify-self-center text-xl")], [
      h.text("Custom game selected"),
    ]),
    h.div(
      [class("grid grid-cols-2 justify-items-center lg:px-50 md:px-25 sm:px-5")],
      [
        form_field.render(
          form_field.new()
            |> form_field.float_label(form_field.Auto),
          [],
          [
            h.label([attribute("slot", "label"), for("start")], [
              h.text("Start"),
            ]),
            h.input([
              event.on_input(CustomStartChanged),
              id("start"),
              required(True),
            ]),
            h.span([attribute("slot", "hint")], [
              h.text("Initial Wikipedia topic"),
            ]),
          ],
        ),
        form_field.render(
          form_field.new() |> form_field.float_label(form_field.Auto),
          [],
          [
            h.label([attribute("slot", "label"), for("goal")], [h.text("Goal")]),
            h.input([
              event.on_input(CustomGoalChanged),
              id("goal"),
              required(True),
            ]),
            h.span([attribute("slot", "hint")], [
              h.text("Goal/target Wikipedia topic"),
            ]),
          ],
        ),
        ..error_line
      ],
    ),
    h.div([class("justify-self-center")], [
      button.render(
        button.new("Continue", button.Filled)
          |> button.disabled(or(string.is_empty(st), string.is_empty(gl))),
        [
          event.on_click(CustomEndPointsSelected),
        ],
      ),
    ]),
  ])
}

/// custom_error_line is a helper function that generates the error messages
/// for the custom game setup screen.
fn custom_error_line(
  start_error: Option(String),
  goal_error: Option(String),
) -> #(List(Element(Msg)), String) {
  case start_error, goal_error {
    Some(se), Some(ge) -> #(
      [
        h.div([class("bg-error self-center justify-self-start")], [h.text(se)]),
        h.div([class("bg-error self-center justify-self-start")], [h.text(ge)]),
      ],
      "grid-rows-4",
    )
    Some(se), None -> #(
      [h.div([class("bg-error self-center justify-self-start")], [h.text(se)])],
      "grid-rows-4",
    )
    None, Some(ge) -> #(
      [
        h.div([class("bg-error self-center justify-self-start col-2")], [
          h.text(ge),
        ]),
      ],
      "grid-rows-4",
    )
    None, None -> #([element.none()], "grid-rows-3")
  }
}

/// error_message formats an rsvp error
///
fn error_message(rsvp_error: Option(String)) -> Element(Msg) {
  case rsvp_error {
    Some(message) ->
      h.div([class("bg-error self-center justify-self-start")], [
        h.text(message),
      ])

    None -> element.none()
  }
}

/// goal builds an HTML element that congratulates the user
/// when the Goal topic has been navigated
///
/// Parameters:
///   state: current game state (one of the setup phase values)
///
fn goal(state: State) -> Element(Msg) {
  case state {
    Completed ->
      h.div([class("grid font-bold place-content-center text-xl bg-teal-500")], [
        h.div([class("justify-self-center")], [h.text("Goal! Goal! Goal!")]),
      ])
    _ -> element.none()
  }
}

/// random builds the section of the setup page which handles
/// randomly generated Start and Goal topics
///
/// Parameters:
///   start: the selected starting Wiki topic
///   goal: the selected Wiki goal topic
///   rsvp_error: error interacting with or from the API server
///
fn random(
  start: EP(Start),
  goal: EP(Goal),
  rsvp_error: Option(String),
) -> Element(Msg) {
  let EP(gl): EP(Goal) = goal
  let EP(st): EP(Start) = start
  let loading = case st, gl, rsvp_error {
    "", "", None ->
      h.div([class("justify-self-center")], [
        h.i(
          [
            class(
              "fa-solid fa-spinner fa-spin-pulse justify-self-center text-5xl",
            ),
          ],
          [element.none()],
        ),
      ])
    "", "", Some(_) -> error_message(rsvp_error)
    _, _, _ -> element.none()
  }
  h.div([class("grid grid-rows-[1fr_1fr_2fr_1fr_1fr]")], [
    h.div([class("self-center justify-self-center text-xl")], [
      h.text("Random game selected"),
    ]),
    h.div(
      [
        class(
          "grid grid-cols-2 lg:px-50 md:px-25 sm:px-5 gap-1 content-center justify-items-center",
        ),
      ],
      [
        h.span([class("font-bold")], [h.text("Start: ")]),
        h.span([class("font-bold")], [h.text("Goal: ")]),
      ],
    ),
    loading,
    h.div(
      [
        class(
          "grid grid-cols-2 lg:px-50 md:px-25 sm:px-5 gap-1 justify-items-center",
        ),
      ],
      [h.p([], [h.text(st)]), h.p([], [h.text(gl)])],
    ),
    h.div(
      [
        class(
          "grid grid-cols-2 lg:px-50 md:px-25 sm:px-5 gap-1 justify-items-center",
        ),
      ],
      [
        button.render(button.new("Continue", button.Filled), [
          class("self-center"),
          event.on_click(RandomEndPointsDisplayed),
        ]),
        button.render(button.new("Deal again", button.Filled), [
          class("self-center"),
          event.on_click(RedrawRandom),
        ]),
      ],
    ),
  ])
}

//
// -----------------------------------------------------------------------------
//
// GAME PLAY PAGE and HELPERS------------------------------------------------------------
//
// -----------------------------------------------------------------------------

/// playing builds the game play page
///
/// Parameters:
///   model: the current Model
///
fn playing(model: Model) -> List(Element(Msg)) {
  let EP(g): EP(Goal) = active_goal(model.endpoints)
  let target =
    h.div([class("grid grid-cols-2 gap-4 text-xl")], [
      h.p([class("justify-self-end self-center")], [h.text("Target:")]),
      h.span([class("font-bold justify-self-start self-center")], [h.text(g)]),
    ])
  [
    h.div([class("grid gap-4")], [
      target,
      playing_controls(model.state, model.navigation),
      goal(model.state),
      progress(model.steps, model.elapsed),
    ]),
    wiki(model.wiki_html, model.pending, model.rsvp_error),
  ]
}

/// playing_controls builds a button bar for navigation and game control
///
/// Parameters:
///   state: current game state (one of the setup phase values)
///   nav: the Navigation state
///
fn playing_controls(state: State, nav: Navigation) -> Element(Msg) {
  let #(back, fwd) = navigation_possible(nav)
  let back_disablement = case state, back {
    Completed, _ | _, False -> disabled(True)
    _, True -> disabled(False)
  }
  let fwd_disablement = case state, fwd {
    Completed, _ | _, False -> disabled(True)
    _, True -> disabled(False)
  }
  let cols = case state {
    Playing | Completed | ReadyToPlay -> class("grid-cols-3")
    Paused -> class("grid-cols-5")
    _ -> class("grid-cols-3")
  }
  h.div([class("grid lg:px-50 md:px-25 sm:px-5 gap-1 justify-center"), cols], [
    button.render(button.new("Back", button.Filled), [
      event.on_click(NavigateBack),
      back_disablement,
    ]),
    button.render(button.new("Forward", button.Filled), [
      event.on_click(NavigateForward),
      fwd_disablement,
    ]),
    case state {
      Playing ->
        button.render(button.new("Pause", button.Filled), [
          event.on_click(GamePaused),
        ])
      Paused ->
        button.render(button.new("Continue", button.Filled), [
          event.on_click(GameResumed),
        ])
      ReadyToPlay ->
        button.render(button.new("Play", button.Filled), [
          event.on_click(GameStarted),
        ])
      _ -> element.none()
    },
    case state {
      Completed | Paused ->
        button.render(button.new("New Game", button.Filled), [
          event.on_click(NewGame),
        ])
      _ -> element.none()
    },
    case state {
      Paused ->
        button.render(button.new("Restart", button.Filled), [
          event.on_click(RestartGame),
        ])
      _ -> element.none()
    },
  ])
}

/// format_seconds accepts an integer time in seconds and returns a
/// string in HH:MM:SS format
///
fn format_seconds(total_seconds: Int) -> String {
  let hours = total_seconds / 3600
  let minutes = { total_seconds % 3600 } / 60
  let seconds = total_seconds % 60

  let hh = case hours {
    h if h < 10 -> "0" <> int.to_string(h)
    _ -> int.to_string(hours)
  }

  let mm = case minutes {
    m if m < 10 -> "0" <> int.to_string(m)
    _ -> int.to_string(minutes)
  }

  let ss = case seconds {
    s if s < 10 -> "0" <> int.to_string(s)
    _ -> int.to_string(seconds)
  }

  hh <> ":" <> mm <> ":" <> ss
}

/// progress builds an HTML element that reports the number of
/// steos taken to date and the elapsed time for game play
///
/// Parameters:
///   steps: steps taken so far
///   elapsed: number of seconds of active gane play so far
///
fn progress(steps: Int, elapsed: Int) -> Element(Msg) {
  h.div([class("grid grid-cols-2 gap-4")], [
    h.div([class("grid content-center justify-end text-xl")], [
      h.p([], [
        h.text("Steps: "),
        h.span([class("font-bold")], [h.text(int.to_string(steps))]),
      ]),
    ]),
    h.div([class("grid content-center justify-start text-xl")], [
      h.p([], [
        h.text(" Elapsed: "),
        h.span([class("font-bold")], [h.text(format_seconds(elapsed))]),
      ]),
    ]),
  ])
}

/// special_on_click returns a Lustre Event (as an Attribute)
/// which receives Click events from anywhere within the
/// encapsulted Wiki page that is currently displayed. It discards
/// any click that is not on an A element. For clicks on A
/// elements it returns the HREF from the A element
///
/// Parameters:
///   msg_constructor: the Msg function thaty receives either the
///                    HREF or an error
///
fn special_on_click(msg_constructor: fn(String) -> msg) {
  event.on("click", {
    use tag <- decode.subfield(["target", "tagName"], decode.string)

    case tag {
      "A" ->
        decode.at(["target", "href"], decode.string)
        |> decode.map(msg_constructor)
      _ -> decode.failure(msg_constructor(""), "<a>")
    }
  })
}

/// wiki builds the HTML which encapsultes the <body> of a Wiki page
///
/// Parameters:
///   page: content of the Wiki page between the <body> and </body> tags
///   pending: when not empty this indicates that a page fetch Effect is
///            in flight and so a spinner should be shown
///   rsvp_error: error interacting with or from the API server
fn wiki(
  page: String,
  pending: String,
  rsvp_error: Option(String),
) -> Element(Msg) {
  let loading = case pending, rsvp_error {
    "", None -> element.none()
    "", Some(e) -> h.div([class("bg-error justify-self-center")], [h.text(e)])
    _, _ ->
      h.div([class("justify-self-center")], [
        h.i(
          [
            class(
              "fa-solid fa-spinner fa-spin-pulse justify-self-center text-5xl",
            ),
          ],
          [element.none()],
        ),
      ])
  }
  h.div([class("grid grid-rows-2")], [
    loading,
    // Using `unsafe_raw_html` is necessary here to render the HTML content
    // from the Wikipedia API. This is safe because the content is from a
    // trusted source. If the content were from an untrusted source, it would
    // need to be sanitized to prevent XSS attacks.
    element.unsafe_raw_html(
      "",
      "div",
      [
        id("wiki"),
        special_on_click(Click)
          |> event.prevent_default()
          |> event.stop_propagation(),
      ],
      page,
    ),
  ])
}
