package Helper

type Array_string struct {
	Values []string
}
func (a *Array_string) Append(value string) {
	a.Values = append(a.Values, value)
}
