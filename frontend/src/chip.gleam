//// Moduke chip contains functions that View uses to generate an HTML chip

import gleam/list

import lustre/attribute.{type Attribute, class}
import lustre/element.{type Element}
import lustre/element/html.{div}

import colours.{
  type Colour, Primary, PrimaryContainer, Secondary, SecondaryContainer,
  Tertiary, TertiaryContainer,
}
import size.{type Size, Icon, Large, Medium, Small}

// Variant represents the different flavours of chip that this
// nodule can generate
//
pub type Variant {
  Flat
  Ghost
  Light
  Solid
  Outlined
}

/// chip is the interface to this module. It returns a Lustre
/// Element representing an HTML chip in the requested flavour
///
/// Parameters:
///   variety: the flavour of chip desired
///   size: the chip size
///   colour: the Material color
///   attributes: additional attributes (e.g. layout)
///   children: child Elements
///
pub fn chip(
  variety: Variant,
  size: Size,
  colour: Colour,
  attributes: List(Attribute(a)),
  children: List(Element(a)),
) -> Element(a) {
  let attr = [
    class("gap-2 rounded-full transition-all"),
    size_classes(size),
    colour_classes(colour, variety),
    ..attributes
  ]
  div(list.append(variant_classes(variety), attr), children)
}

/// colour_classes generates Tailwind CSS classes for Material colors
///
fn colour_classes(group: Colour, variety: Variant) -> Attribute(a) {
  case group, variety {
    Primary, Solid -> class("bg-primary text-primary")
    Primary, Outlined -> class("outline-primary text-primary")
    Primary, Light -> class("text-primary")
    Primary, Flat -> class("bg-primary/20 text-primary")
    PrimaryContainer, Solid ->
      class("bg-primary-container text-on-primary-container")
    PrimaryContainer, Outlined ->
      class("outline-on-primary-container text-on-primary")
    PrimaryContainer, Light -> class("text-on-primary-container")
    PrimaryContainer, Flat ->
      class("bg-primary-container/20 text-on-primary-container")
    Secondary, Solid -> class("bg-secondary text-secondary")
    Secondary, Outlined -> class("outline-secondary text-secondary")
    Secondary, Light -> class("text-secondary")
    Secondary, Flat -> class("bg-secondary/20 text-secondary")
    SecondaryContainer, Solid ->
      class("bg-secondary-container text-on-secondary-container")
    SecondaryContainer, Outlined ->
      class("outline-secondaryy-on-container text-secondaryy-on-container")
    SecondaryContainer, Light -> class("text-on-secondary-container")
    SecondaryContainer, Flat ->
      class("bg-secondary-container/20 text-on-secondary-container")
    Tertiary, Solid -> class("bg-tertiary text-tertiary")
    Tertiary, Outlined -> class("outline-tertiary text-tertiary")
    Tertiary, Light -> class("text-tertiary")
    Tertiary, Flat -> class("bg-tertiary/20 text-tertiary")
    TertiaryContainer, Solid ->
      class("bg-tertiary-container text-on-tertiary-container")
    TertiaryContainer, Outlined ->
      class("outline-on-tertiary-container text-on-tertiary-container")
    TertiaryContainer, Light -> class("text-on-tertiary-container")
    TertiaryContainer, Flat ->
      class("bg-tertiary-container/20 text-on-tertiary-container")
    _, _ -> class("bg-primary text-primary")
  }
}

/// size_classes generates Tailwind CSS classes for chip sizing
///
fn size_classes(size: Size) -> Attribute(a) {
  case size {
    Small -> class("px-2 py-1 text-xs")
    Medium -> class("px-3 py-1 text-sm")
    Large -> class("px-4 py-1 text-base")
    Icon -> class("p-1")
  }
}

/// variant_classes generates Tailwind CSS classes for chip styling
///
fn variant_classes(variety: Variant) -> List(Attribute(a)) {
  case variety {
    Flat | Light | Solid -> []
    Outlined -> [class("outline outline-1")]
    Ghost -> [class("bg-transparent border-2")]
  }
}
