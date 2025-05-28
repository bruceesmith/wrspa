//// Module Msg defines the Msg type used by the Lustre Model-View-Update loop

import rsvp

import endpoints.{type EP, type Goal, type Start}

pub type Msg {
  // Setup messages
  CustomEndPointsSelected
  CustomGoalChanged(String)
  CustomSelected
  CustomStartChanged(String)
  DarkModeFetched(Bool)
  RandomEndPointsDisplayed
  RandomSelected
  SpecialRandomFetched(Result(#(EP(Goal), EP(Start)), rsvp.Error))

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
