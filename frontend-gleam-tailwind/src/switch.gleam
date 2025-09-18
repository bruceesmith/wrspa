//// Module switch contains functions to build an HTML toggle switch

import lustre/attribute.{checked, class, for, id, type_}
import lustre/element.{type Element}
import lustre/element/html as h
import lustre/event

/// switch is the interface to this module. It returns a Lustre
/// Element representing an HTML toggle switch
/// 
/// Parameters:
///   set: initial state for the switch
///   label: label for the switch
///   on_check: message sent when the state of the switch is changed by the user clicking on iy
///
pub fn switch(set: Bool, label: String, on_check: fn(Bool) -> a) -> Element(a) {
  h.div([class("grid grid-rows-1 grid-cols-2")], [
    h.text(label),
    h.div(
      [
        class(
          "relative inline-block w-[52px] h-[32px] justify-self-center self-center",
        ),
      ],
      [
        h.input([
          checked(set),
          id("switch-component"),
          type_("checkbox"),
          class(
            "peer appearance-none w-[52px] h-[32px] bg-surface-container-highest rounded-full checked:bg-primary cursor-pointer transition-colors duration-300",
          ),
          event.on_check(on_check),
        ]),
        h.label(
          [
            for("switch-component"),
            class(
              "absolute top-[4px] left-[2px] w-[24px] h-[24px] bg-outline peer-checked:bg-on-primary-container rounded-full border border-slate-300 shadow-sm transition-transform duration-300 peer-checked:translate-x-[25px] peer-checked:border-slate-800 cursor-pointer",
            ),
          ],
          [],
        ),
      ],
    ),
  ])
}
