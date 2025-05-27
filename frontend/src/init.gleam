//// Module init contains the functions for SPA initialisation and
//// game reset

import gleam/option.{None}

import lustre/effect.{type Effect}

import effects.{dark_mode_on, special_random}
import model.{type Model, ChoosingGame, Endpoints, Model}
import msg.{type Msg, DarkModeFetched, SpecialRandomFetched}
import navigation

/// init is called by Model-View-Update at application initialisation. It establishes
/// the initial Model record
///
pub fn init(_args) -> #(Model, Effect(Msg)) {
  #(
    initial(),
    effect.batch([
      special_random(SpecialRandomFetched),
      dark_mode_on(DarkModeFetched),
    ]),
  )
}

/// initial performs the build of a clean Model. It is split out from
/// init so that it can be called when the used requests to start a new game
///
pub fn initial() -> Model {
  Model(
    dark: False,
    displayed: "",
    elapsed: 0,
    endpoints: Endpoints("", "", "", "", "", ""),
    goal_error: None,
    navigation: navigation.new(),
    pending: "",
    rsvp_error: None,
    start_error: None,
    state: ChoosingGame,
    steps: 0,
    timer_id: 0,
    wiki_html: "",
  )
}

/// reset is called when the user requests to start a particular game
/// over again. It differs from initial in that the existing start and
/// goal are preserved
///
pub fn reset(model: Model) -> Model {
  Model(
    ..model,
    displayed: "",
    elapsed: 0,
    goal_error: None,
    navigation: navigation.new(),
    pending: model.endpoints.actual_start,
    rsvp_error: None,
    start_error: None,
    state: model.ReadyToPlay,
    steps: 0,
    timer_id: 0,
    wiki_html: "",
  )
}
