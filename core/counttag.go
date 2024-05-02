package core

type CountTag struct {
	Tag   string
	Count int
}

type Tags []CountTag

func (t Tags) Tags() []string {
	var result = make([]string, len(t))
	for i, _ := range t {
		result[i] = t[i].Tag
	}
	return result
}
