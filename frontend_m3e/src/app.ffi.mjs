/*
This JS module contains small functions that are called inside
Lustre Effect processing
*/

// clearInterval()
export function clear_interval(id) {
  window.clearInterval(id);
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
