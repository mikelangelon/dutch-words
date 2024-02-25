package core

type FormData struct {
	Errors map[string]string
}

func NewFormData() FormData {
	return FormData{
		Errors: make(map[string]string),
	}
}
