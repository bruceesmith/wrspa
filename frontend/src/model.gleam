//// Module model defines the Model record that is central to the Lustre
//// Model-View-Update architecture.

import gleam/option.{type Option}

import navigation.{type Navigation}

/// Model is the central data structure of this SPA
///
pub type Model {
  Model(
    dark: Bool,
    displayed: String,
    elapsed: Int,
    endpoints: Endpoints,
    goal_error: Option(String),
    navigation: Navigation,
    pending: String,
    rsvp_error: Option(String),
    start_error: Option(String),
    state: State,
    steps: Int,
    timer_id: Int,
    wiki_html: String,
  )
}

/// Endpoints defines the start and goal Wiki topics
///
pub type Endpoints {
  Endpoints(goal: String, start: String)
}

/// State defines the current state of game setup and
/// play. It is the key value which drives the View
/// function's display of HTML
///
pub type State {
  ChoosingGame
  Completed
  CustomGame
  Playing
  Paused
  RandomGame
  ReadyToPlay
}
