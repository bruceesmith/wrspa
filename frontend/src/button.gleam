//// Moduke button contains functions that View uses to generate an HTML button

import gleam/list

import lustre/attribute.{type Attribute, class}
import lustre/element.{type Element}
import lustre/element/html

import colours.{
  type Colour, Primary, PrimaryContainer, Secondary, SecondaryContainer,
  Tertiary, TertiaryContainer,
}
import size.{type Size, Icon, Large, Medium, Small}

// Variant represents the different flavours of button that this
// nodule can generate
//
pub type Variant {
  Flat
  Ghost
  Light
  Outlined
  Solid
}

/// button is the interface to this module. It returns a Lustre
/// Element representing an HTML button in the requested flavour
///
/// Parameters:
///   variety: the flavour of button desired
///   size: the button size
///   colour: the Material color
///   attributes: additional attributes (e.g. layout)
///   children: child Elements, e.g an icon
///
pub fn button(
  variety: Variant,
  size: Size,
  colour: Colour,
  attributes: List(Attribute(a)),
  children: List(Element(a)),
) -> Element(a) {
  let attr = [
    class(
      "justify-self-center transition-all active:enabled:scale-[98%] disabled:opacity-50 disabled:cursor-not-allowed",
    ),
    size_classes(size),
    colour_classes(colour),
    ..attributes
  ]
  html.button(list.append(variant_classes(variety), attr), children)
}

/// colour_classes generates Tailwind CSS classes for Material colors
///
fn colour_classes(group: Colour) -> Attribute(a) {
  case group {
    Primary -> class("bg-primary text-primary hover:text-primary/50")
    PrimaryContainer ->
      class(
        "bg-primary-container text-on-primary-container hover:text-on-primary-container/50",
      )
    Secondary -> class("bg-secondary text-secondary hover:text-secondary/50")
    SecondaryContainer ->
      class(
        "bg-secondary-container text-on-secondary-container hover:text-on-secondary-container/50",
      )
    Tertiary -> class("bg-tertiary text-tertiary hover:text-tertiary/50")
    TertiaryContainer ->
      class(
        "bg-tertiary-container text-on-tertiary-container hover:text-on-tertiary-container/50",
      )
  }
}

/// size_classes generates Tailwind CSS classes for button sizing
///
fn size_classes(size: Size) -> Attribute(a) {
  case size {
    Small -> class("rounded-sm px-3.5 py-1.5 text-sm")
    Medium -> class("rounded-md px-4 py-2 text-base")
    Large -> class("rounded-lg px-5 py-2.5 text-lg")
    Icon -> class("rounded-md p-2")
  }
}

/// variant_classes generates Tailwind CSS classes for buttonm styling
///
fn variant_classes(variety: Variant) -> List(Attribute(a)) {
  case variety {
    Solid | Flat -> []
    Outlined | Ghost -> [class("bg-transparent border-2")]
    Light -> [class("bg-transparent")]
  }
}
