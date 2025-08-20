//// Module update provides the update() function, one of the
//// Lustre Model-View-Update triumvirate, and a series of helper
//// functions

import gleam/int

// import gleam/io
import gleam/option.{None, Some}
import gleam/regexp
import gleam/string
import gleam/uri

import lustre/effect.{type Effect}
import rsvp

import effects.{
  dark_mode_on, get_wiki_page, scroll_up, special_random, start_timer,
  stop_timer,
}
import endpoints.{
  type EP, type Goal, type Start, EP, active_from_custom, active_from_random,
  active_goal, active_start, create_endpoint, custom_start, new_random,
  random_start, update_custom_goal, update_custom_start, update_random,
}
import init.{initial, reset}
import model.{
  type Model, Completed, CustomGame, Model, Paused, Playing, RandomGame,
  ReadyToPlay,
}
import msg.{
  type Msg, Click, CustomEndPointsSelected, CustomGoalChanged, CustomSelected,
  CustomStartChanged, DarkModeFetched, DarkModeSetting, GamePaused, GameResumed,
  GameStarted, NavigateBack, NavigateForward, NewGame, RandomEndPointsDisplayed,
  RandomSelected, RedrawRandom, RestartGame, Scrolled, SpecialRandomFetched,
  TimerReturnedID, TimerStopped, TimerTick, WikiPageFetched,
}
import navigation.{navigate_back, navigate_forward, navigated_to}

// -----------------------------------------------------------------------------
//
// UPDATE ----------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

/// update is called in the Update phase of Model-View-Update when a message is
/// received, either when a user action in the displayed HTML triggers a Msg, or
/// when an Effect triggers a Msg
///
/// Parameters:
///   model: the current SPA state, an instance of the Model record
///   msg: the Msg triggered either by user action or by an Effect
///
pub fn update(model: Model, msg: Msg) -> #(Model, Effect(Msg)) {
  case msg {
    // Setup messages
    DarkModeFetched(set) -> #(Model(..model, dark: set), effect.none())

    CustomSelected -> #(Model(..model, state: CustomGame), effect.none())

    CustomEndPointsSelected -> {
      let subject = {
        let EP(st): EP(Start) = model.endpoints |> custom_start
        "/wiki/" <> st
      }
      #(
        Model(
          ..model,
          endpoints: model.endpoints |> active_from_custom,
          pending: subject,
          state: ReadyToPlay,
        ),
        get_wiki_page(subject, WikiPageFetched),
      )
    }

    CustomStartChanged(val) -> {
      case check_subject(val) {
        Ok(Nil) -> update_custom_endpoint(model, val, True)
        Error(e) -> {
          #(Model(..model, start_error: Some(e)), effect.none())
        }
      }
    }

    CustomGoalChanged(val) -> {
      case check_subject(val) {
        Ok(Nil) -> update_custom_endpoint(model, val, False)
        Error(e) -> {
          #(Model(..model, goal_error: Some(e)), effect.none())
        }
      }
    }

    RandomSelected -> #(Model(..model, state: RandomGame), effect.none())

    RandomEndPointsDisplayed -> {
      let EP(st): EP(Start) = model.endpoints |> random_start
      let subject = "/wiki/" <> st
      #(
        Model(
          ..model,
          endpoints: model.endpoints |> active_from_random,
          pending: subject,
          state: ReadyToPlay,
        ),
        get_wiki_page(subject, WikiPageFetched),
      )
    }

    SpecialRandomFetched(Ok(sr)) -> #(
      Model(
        ..model,
        endpoints: model.endpoints |> update_random(sr.0, sr.1),
        rsvp_error: None,
      ),
      effect.none(),
    )

    SpecialRandomFetched(Error(err)) -> #(
      Model(..model, rsvp_error: Some(rsvp_error_to_string(err))),
      effect.none(),
    )

    // Game play messages
    Click(href) -> check_click(model, href)

    GamePaused -> #(
      Model(..model, state: Paused),
      stop_timer(model.timer_id, TimerStopped),
    )

    GameResumed -> #(
      Model(..model, state: Playing),
      start_timer(1000, TimerReturnedID, TimerTick),
    )

    GameStarted -> #(
      Model(..model, state: Playing),
      start_timer(1000, TimerReturnedID, TimerTick),
    )

    NewGame -> #(
      initial(),
      effect.batch([
        special_random(SpecialRandomFetched),
        dark_mode_on(DarkModeFetched),
      ]),
    )

    RedrawRandom -> #(
      Model(..model, endpoints: model.endpoints |> new_random),
      special_random(SpecialRandomFetched),
    )

    RestartGame -> {
      let EP(st): EP(Start) = model.endpoints |> active_start
      #(reset(model), get_wiki_page("/wiki/" <> st, WikiPageFetched))
    }

    Scrolled -> #(model, effect.none())

    WikiPageFetched(Ok(response)) -> {
      let EP(g): EP(Goal) =
        model.endpoints
        |> active_goal
      let #(st, effie) = case
        string.lowercase(model.pending) == "/wiki/" <> string.lowercase(g)
      {
        True -> #(
          Completed,
          effect.batch([
            stop_timer(model.timer_id, TimerStopped),
            scroll_up(Scrolled),
          ]),
        )
        False -> #(model.state, scroll_up(Scrolled))
      }
      {
        let dest = model.pending
        #(
          Model(
            ..model,
            rsvp_error: None,
            navigation: navigated_to(dest, model.navigation),
            pending: "",
            state: st,
            wiki_html: response,
          ),
          effie,
        )
      }
    }

    WikiPageFetched(Error(err)) -> #(
      Model(
        ..model,
        wiki_html: "failed to fetch",
        rsvp_error: Some(rsvp_error_to_string(err)),
      ),
      effect.none(),
    )

    // Navigation messages
    NavigateBack -> {
      case model.state {
        Playing -> {
          let #(navigation, destination) = navigate_back(model.navigation)
          case destination {
            "" -> #(model, effect.none())
            _ -> #(
              Model(
                ..model,
                navigation: navigation,
                pending: destination,
                steps: model.steps + 1,
              ),
              get_wiki_page(destination, WikiPageFetched),
            )
          }
        }
        _ -> #(model, effect.none())
      }
    }

    NavigateForward -> {
      case model.state {
        Playing -> {
          let #(navigation, destination) = navigate_forward(model.navigation)
          case destination {
            "" -> #(model, effect.none())
            _ -> #(
              Model(
                ..model,
                navigation: navigation,
                pending: destination,
                steps: model.steps + 1,
              ),
              get_wiki_page(destination, WikiPageFetched),
            )
          }
        }
        _ -> #(model, effect.none())
      }
    }

    // Timer messages
    TimerReturnedID(id) -> {
      #(Model(..model, timer_id: id), effect.none())
    }

    TimerStopped -> {
      #(model, effect.none())
    }

    TimerTick -> {
      let elapsed = model.elapsed + 1
      #(Model(..model, elapsed:), effect.none())
    }

    // Settings messages
    DarkModeSetting(mode) -> #(Model(..model, dark: mode), effect.none())
  }
}

// -----------------------------------------------------------------------------
//
// HELPERS----------------------------------------------------------------------
//
// -----------------------------------------------------------------------------

/// check_click is called during game play when the user clicks on an A element inside
/// the encapsulated Wiki page. If the HREF begins with '/wiki/' then the API server is
/// called (as an Effect) to fetch that Wiki topic, otherwise the click is silently ignored
///
/// Parameters:
///   model: the current SPA state, an instance of the Model record
///   href: the HREF from the A element
///
fn check_click(model: Model, href: String) -> #(Model, Effect(Msg)) {
  case model.state {
    Playing -> {
      let maybe_url = uri.parse(href)
      case maybe_url {
        Ok(url) ->
          case string.starts_with(url.path, "/wiki/") {
            True -> #(
              Model(..model, pending: url.path, steps: model.steps + 1),
              get_wiki_page(url.path, WikiPageFetched),
            )
            False -> #(model, effect.none())
          }
        Error(_) -> #(model, effect.none())
      }
    }
    _ -> #(model, effect.none())
  }
}

/// check_subject is called during game setup when the user types inside the
/// 'start' or 'goal' input fields. It validates that only legitimate characters
/// have been typed, and that the field is not empty
///
/// Parameters:
///   subject: the current field value
///
fn check_subject(subject: String) -> Result(Nil, String) {
  let re_result = regexp.from_string("^[a-zA-Z0-9-._~()]+$")
  case re_result {
    Ok(re) ->
      case string.is_empty(subject) {
        False ->
          case regexp.check(re, subject) {
            True -> Ok(Nil)
            False -> Error("only valid characters are a-zA-Z0-9-._~")
          }
        True -> Error("field cannot be empty")
      }
    Error(err) -> Error(err.error)
  }
}

/// rsvp_error_to_string translates an error from the rsvp module to a string
///
fn rsvp_error_to_string(err: rsvp.Error) -> String {
  case err {
    rsvp.BadBody -> "The server sent a response with a bad body."

    rsvp.BadUrl(u) -> "The URL " <> u <> " is invalid."

    rsvp.HttpError(resp) ->
      "There was an HTTP error: " <> int.to_string(resp.status)

    rsvp.JsonError(_) ->
      "There was an error decoding the JSON response from the server."

    rsvp.NetworkError ->
      "There was a network error. Please check your internet connection."

    rsvp.UnhandledResponse(_) ->
      "The server sent a response that could not be handled."
  }
}

/// update_custom_endpoint is a helper function that updates either the custom
/// start or goal endpoint.
///
fn update_custom_endpoint(
  model: Model,
  val: String,
  is_start: Bool,
) -> #(Model, Effect(Msg)) {
  case check_subject(val) {
    Ok(Nil) -> {
      let new_endpoints = case is_start {
        True -> model.endpoints |> update_custom_start(create_endpoint(val))
        False -> model.endpoints |> update_custom_goal(create_endpoint(val))
      }
      let new_model =
        Model(
          ..model,
          endpoints: new_endpoints,
          start_error: case is_start {
            True -> None
            False -> model.start_error
          },
          goal_error: case is_start {
            True -> model.goal_error
            False -> None
          },
        )
      #(new_model, effect.none())
    }
    Error(e) -> {
      let new_model = case is_start {
        True -> Model(..model, start_error: Some(e))
        False -> Model(..model, goal_error: Some(e))
      }
      #(new_model, effect.none())
    }
  }
}
