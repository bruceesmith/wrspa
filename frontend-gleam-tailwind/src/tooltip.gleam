//// Module tooltip contains functions for creating tool tip elements

import lustre/attribute.{class}
import lustre/element.{type Element}
import lustre/element/html as h

/// Position determines where the tooltip is positioned relative to
/// the containing element
/// 
pub type Position {
  Bottom
  Left
  Right
  Top
}

/// tooltip returns a tooltip element. For the tooltip to work, this 
/// returned element must be contained inside another element with the
/// 'tooltip' class
/// 
pub fn tooltip(text: String, pos: Position) -> Element(a) {
  let positioning = case pos {
    Bottom -> class("absolute z-1 top-[100%] left-[50%] ml-[-60px]")
    Left -> class("absolute z-1 top-[-5px] right-[105%]")
    Right -> class("absolute z-1 top-[-5px] left-[106%]")
    Top -> class("absolute z-1 bottom-[100%] left-[50%] ml-[-60px]")
  }
  h.span([class("font-[sans-serif] tooltiptext text-base"), positioning], [
    h.text(text),
  ])
}
