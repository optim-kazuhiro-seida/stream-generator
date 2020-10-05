package sample

//go:generate stream-generator -type=Sample
type Sample struct {
	Str string
	Int int
}
