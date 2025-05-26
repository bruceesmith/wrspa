/*
This JS module contains small functions that are called inside
Lustre Effect processing
*/

// clearInterval()
export function clear_interval(id) {
  window.clearInterval(id);
}

// matchMedia
export function dark_mode() {
  const darkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
  if (darkMode) {
    return true;
  } else {
    return false;
  }
}

// scrollTo
export function scroll_to_top() {
  window.scrollTo(0, 0);
}

// setInterval()
export function set_interval(interval, cb) {
  let id = window.setInterval(cb, interval);
  return id;
}

// setTimeout()
export function set_timeout(delay, cb) {
  window.setTimeout(cb, delay);
}
