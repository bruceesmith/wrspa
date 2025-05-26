package wrserver

import "fmt"

func marshalFailure(function string, err error, response any) string {
	return `{"msg": "unable to marshal API response", ` +
		`"function": "` + function + `", ` +
		`"error": "` + err.Error() + `", ` +
		`"response": "` + fmt.Sprintf("%+v", response) + `"}`
}
