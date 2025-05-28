//// Modue endpoints contains the opaque type Endpoints and a series of helper function

/// Endpoints defines the start and goal Wiki topics
///
pub opaque type Endpoints {
  Endpoints(active: EndpointPair, custom: EndpointPair, random: EndpointPair)
}

/// EP uses phantom typiing to add safety to a simple String
/// 
pub type EP(a) {
  EP(String)
}

/// Goal is a flavour of EP
/// 
pub type Goal {
  Goal
}

/// Start is a flavour of EP
/// 
pub type Start {
  Start
}

/// EndpointPair is an internal type representing a pair of endpoints (a
/// goal and a start)
type EndpointPair {
  EndpointPair(goal: EP(Goal), start: EP(Start))
}

/// active returns the active endpoints ... these are set once a choice
/// has been made of either a custom or a random game
/// 
fn active(ep: Endpoints) -> EndpointPair {
  ep.active
}

/// active_from_custom sets the active endpoints that will be used in a
/// game from the set of randomly selected endpoints
/// 
pub fn active_from_custom(ep: Endpoints) -> Endpoints {
  Endpoints(..ep, active: ep.custom)
}

/// active_from_random sets the active endpoints that will be used in a
/// game from the set of customer entered endpoints
/// 
pub fn active_from_random(ep: Endpoints) -> Endpoints {
  Endpoints(..ep, active: ep.random)
}

/// active_goal returns the active goal
/// 
pub fn active_goal(ep: Endpoints) -> EP(Goal) {
  ep |> active |> goal
}

/// active_start returns the active goal
pub fn active_start(ep: Endpoints) -> EP(Start) {
  ep |> active |> start
}

/// custom returns the manually selected end points
/// 
fn custom(ep: Endpoints) -> EndpointPair {
  ep.custom
}

/// custom_goal returns the manually selected goal
/// 
pub fn custom_goal(ep: Endpoints) -> EP(Goal) {
  ep |> custom |> goal
}

/// custom_start returns the manually selected start
/// 
pub fn custom_start(ep: Endpoints) -> EP(Start) {
  ep |> custom |> start
}

/// ep_from_string creates a type-safe EP from a string
pub fn ep_from_string(str: String) -> EP(a) {
  EP(str)
}

/// goal returns the goal from any endpoint pair
/// 
fn goal(epp: EndpointPair) -> EP(Goal) {
  epp.goal
}

/// new creates a new Endpoints record with all blank end points
/// 
pub fn new() -> Endpoints {
  Endpoints(new_epp(), new_epp(), new_epp())
}

/// new_epp creates a new end point pair with the provided values
/// 
fn new_epp() -> EndpointPair {
  let gl: EP(Goal) = EP("")
  let st: EP(Start) = EP("")
  EndpointPair(gl, st)
}

/// new_random sets a new end point pair into the random field
/// 
pub fn new_random(ep: Endpoints) -> Endpoints {
  Endpoints(..ep, random: new_epp())
}

/// random returns the randomly selected end points
/// 
fn random(ep: Endpoints) -> EndpointPair {
  ep.random
}

/// random_goal returns the randomly selected goal
/// 
pub fn random_goal(ep: Endpoints) -> EP(Goal) {
  ep |> random |> goal
}

/// random_start returns the randomly selected start
/// 
pub fn random_start(ep: Endpoints) -> EP(Start) {
  ep |> random |> start
}

/// start returns the atart field of any pair of end points
/// 
fn start(epp: EndpointPair) -> EP(Start) {
  epp.start
}

/// update_custom sets the pair of custom selected end points
/// 
fn update_custom(epp: EndpointPair, ep: Endpoints) -> Endpoints {
  Endpoints(..ep, custom: epp)
}

/// update_custom_goal sets the goal endpoint in the custom pair
/// 
pub fn update_custom_goal(ep: Endpoints, goal: EP(Goal)) -> Endpoints {
  ep
  |> custom
  |> update_goal(goal)
  |> update_custom(ep)
}

/// update_custom_start sets the start endpoint in the custom pair
/// 
pub fn update_custom_start(ep: Endpoints, start: EP(Start)) -> Endpoints {
  ep
  |> custom
  |> update_start(start)
  |> update_custom(ep)
}

/// update_goal updates the goal field of any end point pair
/// 
fn update_goal(epp: EndpointPair, goal: EP(Goal)) -> EndpointPair {
  EndpointPair(..epp, goal: goal)
}

/// update_random sets a new end point pair into the random field
/// 
pub fn update_random(
  ep: Endpoints,
  goal: EP(Goal),
  start: EP(Start),
) -> Endpoints {
  Endpoints(..ep, random: EndpointPair(goal, start))
}

/// update_start updates the start field of any end point pair
fn update_start(epp: EndpointPair, start: EP(Start)) -> EndpointPair {
  EndpointPair(..epp, start: start)
}
