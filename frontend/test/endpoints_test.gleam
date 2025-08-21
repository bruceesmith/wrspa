import gleeunit/should

import endpoints.{
  type EP, type Endpoints, type Goal, type Start, active_from_custom,
  active_from_random, active_goal, active_start, create_endpoint, custom_goal,
  custom_start, new, new_random, random_goal, random_start, update_custom_goal,
  update_custom_start, update_random,
}

// Tests from endpoints_test.gleam

fn sample_endpoints() -> Endpoints {
  let custom_goal = create_endpoint("Custom Goal")
  let custom_start = create_endpoint("Custom Start")
  let random_goal = create_endpoint("Random Goal")
  let random_start = create_endpoint("Random Start")

  new()
  |> update_custom_goal(custom_goal)
  |> update_custom_start(custom_start)
  |> update_random(random_goal, random_start)
}

pub fn new_test() {
  let ep = new()
  let start: EP(Start) = create_endpoint("")
  let goal: EP(Goal) = create_endpoint("")
  custom_goal(ep)
  |> should.equal(goal)
  custom_start(ep)
  |> should.equal(start)
  random_goal(ep)
  |> should.equal(goal)
  random_start(ep)
  |> should.equal(start)
  active_goal(ep)
  |> should.equal(goal)
  active_start(ep)
  |> should.equal(start)
}

pub fn create_endpoint_test() {
  let ep: EP(String) = create_endpoint("test")
  ep
  |> should.equal(create_endpoint("test"))
}

pub fn update_custom_goal_test() {
  let ep = new()
  let goal = create_endpoint("goal")
  let updated_ep = update_custom_goal(ep, goal)
  custom_goal(updated_ep)
  |> should.equal(goal)
}

pub fn update_custom_start_test() {
  let ep = new()
  let start = create_endpoint("start")
  let updated_ep = update_custom_start(ep, start)
  custom_start(updated_ep)
  |> should.equal(start)
}

pub fn update_random_test() {
  let ep = new()
  let goal = create_endpoint("random_goal")
  let start = create_endpoint("random_start")
  let updated_ep = update_random(ep, goal, start)
  random_goal(updated_ep)
  |> should.equal(goal)
  random_start(updated_ep)
  |> should.equal(start)
}

pub fn active_from_custom_test() {
  let ep = sample_endpoints()
  let active_ep = active_from_custom(ep)
  active_goal(active_ep)
  |> should.equal(create_endpoint("Custom Goal"))
  active_start(active_ep)
  |> should.equal(create_endpoint("Custom Start"))
}

pub fn active_from_random_test() {
  let ep = sample_endpoints()
  let active_ep = active_from_random(ep)
  active_goal(active_ep)
  |> should.equal(create_endpoint("Random Goal"))
  active_start(active_ep)
  |> should.equal(create_endpoint("Random Start"))
}

pub fn new_random_test() {
  let ep = sample_endpoints()
  let new_random_ep = new_random(ep)
  random_goal(new_random_ep)
  |> should.equal(create_endpoint(""))
  random_start(new_random_ep)
  |> should.equal(create_endpoint(""))
}

pub fn active_goal_test() {
  let ep = sample_endpoints()
  let active_ep = active_from_custom(ep)
  active_goal(active_ep)
  |> should.equal(create_endpoint("Custom Goal"))
}

pub fn active_start_test() {
  let ep = sample_endpoints()
  let active_ep = active_from_custom(ep)
  active_start(active_ep)
  |> should.equal(create_endpoint("Custom Start"))
}

pub fn custom_goal_test() {
  let ep = sample_endpoints()
  custom_goal(ep)
  |> should.equal(create_endpoint("Custom Goal"))
}

pub fn custom_start_test() {
  let ep = sample_endpoints()
  custom_start(ep)
  |> should.equal(create_endpoint("Custom Start"))
}

pub fn random_goal_test() {
  let ep = sample_endpoints()
  random_goal(ep)
  |> should.equal(create_endpoint("Random Goal"))
}

pub fn random_start_test() {
  let ep = sample_endpoints()
  random_start(ep)
  |> should.equal(create_endpoint("Random Start"))
}
