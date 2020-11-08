package adapter

type Edge interface {
	Attrs(labelvalues ...interface{})
	Value(label string) interface{}
}
