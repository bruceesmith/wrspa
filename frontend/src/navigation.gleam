//// Module navigation contains functions to navigate back and forward
//// during the Play phase of Wiki Racing. These functions are called
//// in the Update function

/// Navigation is the opaque type which holds navigation date
///
pub opaque type Navigation {
  Navigation(left: List(String), right: List(String))
}

/// new creates an empty Navigation record
///
pub fn new() -> Navigation {
  Navigation([], [])
}

/// navigate_back is called when navigation to a previous WIki page
/// is requested. Back navigation is only possible if there are at least
/// 2 elements in the 'right' list. The first (newest) element
/// of the 'right' list is pushed on to the front of the 'left'
/// list; both the first and second elements of the 'right' list
/// are removed, and the (now removed) second element of the
/// 'right' list is returned
///
pub fn navigate_back(nav: Navigation) -> #(Navigation, String) {
  case nav.right {
    [] | [_] -> {
      #(nav, "")
    }
    [first, second, ..remainder] -> {
      #(Navigation(left: [first, ..nav.left], right: remainder), second)
    }
  }
}

/// navigate_forward is called after having previously navigated back to
/// an earlier Wiki page, and then wishing to navigate forward again.
/// Forward navigation is only possible if there is at least one element
/// in the 'left' list. That first element is removed from the 'left' list,
/// and is returned. The 'right' list is untouched.
///
pub fn navigate_forward(nav: Navigation) -> #(Navigation, String) {
  case nav.left {
    [] -> {
      #(nav, "")
    }
    [only] -> {
      #(Navigation(left: [], right: nav.right), only)
    }
    [first, second, ..remainder] -> {
      #(Navigation(left: [second, ..remainder], right: nav.right), first)
    }
  }
}

/// navigation_possible returns a tuple of Booleans indicating if
/// navigation in the backwards and forwards directions is possible.
///
pub fn navigation_possible(nav: Navigation) -> #(Bool, Bool) {
  case nav.right, nav.left {
    [], [] | [_], [] -> #(False, False)
    [_, _, ..], [] -> #(True, False)
    [], [_] | [], [_, _, ..] | [_], [_] | [_], [_, _, ..] -> #(False, True)
    [_, _, ..], [_] | [_, _, ..], [_, _, ..] -> #(True, True)
  }
}

/// navigated_to is called when the user navigates to a new Wiki page.
/// Once the new page is loaded then its topic is pushed to the front
/// of the 'right' list.
///
pub fn navigated_to(new: String, nav: Navigation) -> Navigation {
  Navigation(..nav, right: [new, ..nav.right])
}
