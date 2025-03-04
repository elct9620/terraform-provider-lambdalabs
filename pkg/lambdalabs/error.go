package lambdalabs

type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

func (e *Error) Error() string {
	return e.Message
}
