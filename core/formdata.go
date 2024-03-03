package core

type FormData struct {
	Tags   []string
	Errors map[string]string
}

func NewFormData(tags []string) FormData {
	return FormData{
		Tags:   tags,
		Errors: make(map[string]string),
	}
}
