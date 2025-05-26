//// Module Msg defines the Msg type used by the Lustre Model-View-Update loop

import rsvp

import model.{type Endpoints}

pub type Msg {
  // Setup messages
  CustomEndPointsSelected
  CustomGoalChanged(String)
  CustomSelected
  CustomStartChanged(String)
  DarkModeFetched(Bool)
  RandomEndPointsDisplayed
  RandomSelected
  SpecialRandomFetched(Result(Endpoints, rsvp.Error))

  // Game play messages
  Click(String)
  GamePaused
  GameResumed
  GameStarted
  NewGame
  RedrawRandom
  RestartGame
  Scrolled
  WikiPageFetched(Result(String, rsvp.Error))

  // Navigation messages
  NavigateBack
  NavigateForward

  // Timer messages
  TimerReturnedID(Int)
  TimerStopped
  TimerTick
}
