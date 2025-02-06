package frontend

// body {
//   margin: 0;
//   overflow: hidden;
//   font-family: Roboto, -apple-system, BlinkMacSystemFont, "Segoe UI", Oxygen,
//   Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;

//   color-scheme: dark;
//   color: #E1E2EB;
//   background-color: #121316;
//   --container-color: #44474E;
//   --primary-color: #ADC6FF;
//   --primary-on-color: #002E68;
//   --surface-container-color: #1F1F23;
//   --secondary-container-color: #3F4759;
//   --secondary-on-container-color: #DAE2F9;
// }

type colour string

const (
	color                     = "#E1E2EB"
	backgroundColor           = "#121316" // 006d77
	containerColor            = "#44474E"
	primaryColor              = "#ADC6FF"
	primaryOnColor            = "#002E68"
	surfaceContainerColor     = "#1F1F23"
	secondaryContainerColor   = "#3F4759"
	secondaryOnContainerColor = "#DAE2F9"
)

func (c colour) String() string {
	return string(c)
}

// https://coolors.co/006d77-83c5be-edf6f9-ffddd2-e29578
// <palette>
//   <color name="Caribbean Current" hex="006d77" r="0" g="109" b="119" />
//   <color name="Tiffany Blue" hex="83c5be" r="131" g="197" b="190" />
//   <color name="Alice Blue" hex="edf6f9" r="237" g="246" b="249" />
//   <color name="Pale Dogwood" hex="ffddd2" r="255" g="221" b="210" />
//   <color name="Atomic tangerine" hex="e29578" r="226" g="149" b="120" />
// </palette>
