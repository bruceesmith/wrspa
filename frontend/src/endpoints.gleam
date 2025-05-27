//// Modue endpoints contains the opaque type Endpoints and a series of helper function

/// Endpoints defines the start and goal Wiki topics
///
pub opaque type Endpoints {
  Endpoints(actual: EndpointPair, custom: EndpointPair, random: EndpointPair)
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

/// actual returns the actual endpoints ... these are set once a choice
/// has been made of either a custom or a random game
/// 
fn actual(ep: Endpoints) -> EndpointPair {
  ep.actual
}

/// actual_from_custom sets the actual endpoints that will be used in a
/// game from the set of randomly selected endpoints
/// 
pub fn actual_from_custom(ep: Endpoints) -> Endpoints {
  Endpoints(..ep, actual: ep.custom)
}

/// actual_from_random sets the actual endpoints that will be used in a
/// game from the set of customer entered endpoints
/// 
pub fn actual_from_random(ep: Endpoints) -> Endpoints {
  Endpoints(..ep, actual: ep.random)
}

/// actual_goal returns the active goal
/// 
pub fn actual_goal(ep: Endpoints) -> EP(Goal) {
  ep |> actual |> goal
}

/// actual_start returns the active goal
pub fn actual_start(ep: Endpoints) -> EP(Start) {
  ep |> actual |> start
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

/// goal returns the goal from any endpoint pair
/// 
fn goal(epp: EndpointPair) -> EP(Goal) {
  epp.goal
}

/// new creates a new Endpoints record with all blank end points
/// 
pub fn new() -> Endpoints {
  Endpoints(
    EndpointPair(EP(""), EP("")),
    EndpointPair(EP(""), EP("")),
    EndpointPair(EP(""), EP("")),
  )
}

/// new_epp creates a new end point pair with the provided values
/// 
fn new_epp(goal goal: EP(Goal), start start: EP(Start)) -> EndpointPair {
  EndpointPair(goal, start)
}

/// new_goal updates the goal field of any end point pair
/// 
fn new_goal(epp: EndpointPair, goal: EP(Goal)) -> EndpointPair {
  EndpointPair(..epp, goal: goal)
}

/// new_random sets a new end point pair into the random field
/// 
pub fn new_random(ep: Endpoints, goal: EP(Goal), start: EP(Start)) -> Endpoints {
  Endpoints(..ep, random: new_epp(goal, start))
}

/// new_start updates the start field of any end point pair
fn new_start(epp: EndpointPair, start: EP(Start)) -> EndpointPair {
  EndpointPair(..epp, start: start)
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

/// set_custom_goal sets the goal endpoint in the custom pair
/// 
pub fn set_custom_goal(ep: Endpoints, goal: EP(Goal)) -> Endpoints {
  ep
  |> custom
  |> new_goal(goal)
  |> set_custom(ep)
}

/// set_custom_start sets the start endpoint in the custom pair
/// 
pub fn set_custom_start(ep: Endpoints, start: EP(Start)) -> Endpoints {
  ep
  |> custom
  |> new_start(start)
  |> set_custom(ep)
}

/// set_custom sets the pair of custom selected end points
/// 
fn set_custom(epp: EndpointPair, ep: Endpoints) -> Endpoints {
  Endpoints(..ep, custom: epp)
}

/// start returns the atart field of any pair of end points
/// 
fn start(epp: EndpointPair) -> EP(Start) {
  epp.start
}
