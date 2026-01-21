import endpoints
import gleam/option.{None, Some}
import gleeunit/should
import init.{initial, reset}
import model
import navigation

pub fn initial_test() {
  let model = initial()

  model.dark |> should.be_false()
  model.elapsed |> should.equal(0)
  model.goal_error |> should.equal(None)
  model.pending |> should.equal("")
  model.rsvp_error |> should.equal(None)
  model.start_error |> should.equal(None)
  case model.state {
    model.ChoosingGame -> Nil
    _ -> should.fail()
  }
  model.steps |> should.equal(0)
  model.timer_id |> should.equal(0)
  model.wiki_html |> should.equal("")
}

pub fn reset_test() {
  let start_endpoint = endpoints.create_endpoint("Gleam (programming language)")
  let model =
    model.Model(
      dark: True,
      elapsed: 100,
      endpoints: endpoints.new()
        |> endpoints.update_custom_start(start_endpoint)
        |> endpoints.active_from_custom(),
      goal_error: Some("Some error"),
      navigation: navigation.new(),
      pending: "Some pending",
      rsvp_error: Some("Some rsvp error"),
      start_error: Some("Some start error"),
      state: model.ReadyToPlay,
      steps: 10,
      timer_id: 123,
      wiki_html: "Some html",
    )

  let reset_model = reset(model)

  reset_model.dark |> should.be_true()
  // preserved
  reset_model.elapsed |> should.equal(0)
  reset_model.endpoints |> should.equal(model.endpoints)
  // preserved
  reset_model.goal_error |> should.equal(None)
  case reset_model.pending {
    "Gleam (programming language)" -> Nil
    _ -> should.fail()
  }
  reset_model.rsvp_error |> should.equal(None)
  reset_model.start_error |> should.equal(None)
  case reset_model.state {
    model.ReadyToPlay -> Nil
    _ -> should.fail()
  }
  reset_model.steps |> should.equal(0)
  reset_model.timer_id |> should.equal(0)
  reset_model.wiki_html |> should.equal("")
}
