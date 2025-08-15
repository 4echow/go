package forms

// error type for validation error messages for forms. Form field -> key in map
type errors map[string][]string

// add error messages for given field to map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

//retrieve the first error message for given field from map
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}