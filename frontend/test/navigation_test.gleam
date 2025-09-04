import gleeunit/should
import navigation.{Navigation, navigate_back, navigate_forward, navigated_to, navigation_possible, new}

pub fn new_test() {
  let nav = new()
  nav
  |> should.equal(Navigation([], []))
}

pub fn navigated_to_test() {
  let nav = new()
  let nav = navigated_to("a", nav)
  let nav = navigated_to("b", nav)

  nav
  |> should.equal(Navigation(left: [], right: ["b", "a"]))
}

pub fn navigate_back_test() {
  let nav = new()
  let nav = navigated_to("a", nav)
  let nav = navigated_to("b", nav)
  let nav = navigated_to("c", nav)

  let #(new_nav, topic) = navigate_back(nav)
  topic
  |> should.equal("b")
  new_nav
  |> should.equal(Navigation(left: ["c"], right: ["a"]))
}

pub fn navigate_back_not_possible_test() {
  let nav = new()
  let nav = navigated_to("a", nav)

  let #(new_nav, topic) = navigate_back(nav)
  topic
  |> should.equal("")
  new_nav
  |> should.equal(Navigation(left: [], right: ["a"]))
}

pub fn navigate_forward_test() {
  let nav = new()
  let nav = navigated_to("a", nav)
  let nav = navigated_to("b", nav)
  let nav = navigated_to("c", nav)

  let #(nav_after_back, _) = navigate_back(nav)
  let #(new_nav, topic) = navigate_forward(nav_after_back)
  topic
  |> should.equal("c")
  new_nav
  |> should.equal(Navigation(left: [], right: ["a"]))
}

pub fn navigate_forward_not_possible_test() {
  let nav = new()
  let nav = navigated_to("a", nav)

  let #(new_nav, topic) = navigate_forward(nav)
  topic
  |> should.equal("")
  new_nav
  |> should.equal(Navigation(left: [], right: ["a"]))
}

pub fn navigation_possible_test() {
  let nav = new()
  navigation_possible(nav)
  |> should.equal(#(False, False))

  let nav = navigated_to("a", nav)
  navigation_possible(nav)
  |> should.equal(#(False, False))

  let nav = navigated_to("b", nav)
  navigation_possible(nav)
  |> should.equal(#(True, False))

  let #(nav, _) = navigate_back(nav)
  navigation_possible(nav)
  |> should.equal(#(False, True))

  let nav = navigated_to("c", nav)
  navigation_possible(nav)
  |> should.equal(#(False, True))
}
