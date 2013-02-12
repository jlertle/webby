package lang

func init() {
	Langs.Add("EnglishGB", Lang{
		"ErrorTemplate":             "<p>{{.}}</p>",
		"ErrMandatory":              "This field is mandatory! (Also first character must not be space or newline)",
		"ErrMinChar":                "Characters must be more than and equal to %d!",
		"ErrMaxChar":                "Characters must be less than and equal to %d!",
		"ErrMustMatchMissing":       "Must Match is missing!",
		"ErrMandatoryCheckbox":      "This checkbox requires marking!",
		"ErrMimeCheck":              "Invalid file mime, valid mime are/is : ",
		"ErrSelectValueMissing":     "Select Value is missing!",
		"ErrSelectOptionIsRequired": "Option is required!",
		"ErrFieldDoesNotExist":      "Field name does not exist!",
		"ErrInvalidEmailAddress":    "Invalid Email Address!",
		"ErrFileRequired":           "File is required!",
		"ErrAntiCSRF":               "Failed CSRF Validation!",
	})
}
