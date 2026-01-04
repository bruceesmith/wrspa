//// Moduke input contains functions that View uses to generate an HTML input element

import gleam/list
import gleam/string

import lustre/attribute.{type Attribute, class}
import lustre/element.{type Element}
import lustre/element/html

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
  Light
  Outlined
  Underlined
}

/// input is the interface to this module. It returns a Lustre
/// Element representing an HTML input element in the requested flavour
///
/// Parameters:
///   variety: the flavour of input element desired
///   size: the field size
///   colour: the Material color
///   attributes: additional attributes (e.g. layout)
///   children: child Elements
///
pub fn input(
  variety: Variant,
  size: Size,
  colour: Colour,
  attributes: List(Attribute(a)),
) -> Element(a) {
  let attr = [
    class("disabled:opacity-50 disabled:cursor-not-allowed transition-all"),
    size_classes(size),
    colour_classes(colour, variety),
    ..attributes
  ]
  html.input(list.append(variant_classes(variety), attr))
}

/// colour_classes generates Tailwind CSS classes for Material colors
///
fn colour_classes(group: Colour, variety: Variant) -> Attribute(a) {
  case group, variety {
    Primary, Flat ->
      class(
        "text-primary placeholder-primary bg-primary/20 hover:enabled:bg-primary/30",
      )
    Primary, Outlined ->
      class(
        "border-primary focus:enabled:border-primary hover:enabled:border-primary text-primary",
      )
    Primary, Light -> class("bg-primary placeholder-primary text-primary")
    Primary, Underlined ->
      class("border-primary placeholder-primary text-primary")
    PrimaryContainer, Flat ->
      class(
        "text-on-primary-container placeholder-on-primary-container bg-primary-container/20 hover:enabled:bg-primary-container/30",
      )
    PrimaryContainer, Outlined ->
      class(
        "border-on-primary-container focus:enabled:border-on-primary-container hover:enabled:border-on-primary-container text-primary",
      )
    PrimaryContainer, Light ->
      class(
        "bg-primary-container placeholder-on-primary-container text-on-primary-container",
      )
    PrimaryContainer, Underlined ->
      class(
        "border-on-primary-container placeholder-on-primary-container text-on-primary-container",
      )
    Secondary, Flat ->
      class(
        "text-secondary placeholder-secondary bg-secondary/20 hover:enabled:bg-secondary/30",
      )
    Secondary, Outlined ->
      class(
        "border-secondary focus:enabled:border-secondary hover:enabled:border-secondary text-secondary",
      )
    Secondary, Light ->
      class("bg-secondary placeholder-secondary text-secondary")
    Secondary, Underlined ->
      class("border-secondary placeholder-secondary text-secondary")
    SecondaryContainer, Flat ->
      class(
        "text-on-secondary-container placeholder-on-secondary-container bg-secondary-container/20 hover:enabled:bg-secondary-container/30",
      )
    SecondaryContainer, Outlined ->
      class(
        "border-on-secondary-container focus:enabled:border-on-secondary-container hover:enabled:border-on-secondary-container text-secondary",
      )
    SecondaryContainer, Light ->
      class(
        "bg-secondary-container placeholder-on-secondary-container text-on-secondary-container",
      )
    SecondaryContainer, Underlined ->
      class(
        "border-on-secondary-container placeholder-on-secondary-container text-on-secondary-container",
      )
    Tertiary, Flat ->
      class(
        "text-tertiary placeholder-tertiary bg-tertiary/20 hover:enabled:bg-tertiary/30",
      )
    Tertiary, Outlined ->
      class(
        "border-tertiary focus:enabled:border-tertiary hover:enabled:border-tertiary text-tertiary",
      )
    Tertiary, Light -> class("bg-tertiary placeholder-tertiary text-tertiary")
    Tertiary, Underlined ->
      class("border-tertiary placeholder-tertiary text-tertiary")
    TertiaryContainer, Flat ->
      class(
        "text-on-tertiary-container placeholder-on-tertiary-container bg-tertiary-container/20 hover:enabled:bg-tertiary-container/30",
      )
    TertiaryContainer, Outlined ->
      class(
        "border-on-tertiary-container focus:enabled:border-on-tertiary-container hover:enabled:border-on-tertiary-container text-tertiary",
      )
    TertiaryContainer, Light ->
      class(
        "bg-tertiary-container placeholder-on-tertiary-container text-on-tertiary-container",
      )
    TertiaryContainer, Underlined ->
      class(
        "border-on-tertiary-container placeholder-on-tertiary-container text-on-tertiary-container",
      )
  }
}

/// size_classes generates Tailwind CSS classes for element sizing
///
fn size_classes(size: Size) -> Attribute(a) {
  case size {
    Small -> class("rounded-sm px-3.5 py-1.5 text-sm")
    Medium -> class("rounded-md px-4 py-2 text-base")
    Large -> class("rounded-lg px-5 py-2.5 text-lg")
    Icon -> class("rounded-sm px-3.5 py-1.5 text-sm")
  }
}

/// variant_classes generates Tailwind CSS classes for element styling
///
fn variant_classes(variety: Variant) -> List(Attribute(a)) {
  case variety {
    Light -> [
      [
        "bg-opacity-0 outline-none",
        "hover:enabled:bg-opacity-20 focus:enabled:bg-opacity-30",
      ]
      |> string.join(" ")
      |> class,
    ]
    Flat -> [class("focus:outline-none focus:ring-2 focus:ring-current")]
    Outlined -> [
      [
        "bg-transparent border-2 border-opacity-50 placeholder-opacity-70",
        "hover:enabled:border-opacity-100 focus:enabled:border-opacity-100",
        "focus:outline-none focus:ring-2 focus:ring-current",
      ]
      |> string.join(" ")
      |> class,
    ]
    Underlined -> [
      [
        "bg-transparent border-opacity-20 border-b-2 outline-none rounded-b-none",
        "focus:enabled:border-opacity-100",
      ]
      |> string.join(" ")
      |> class,
    ]
  }
}
