package encoding

type Coder interface {
	Encode(int64) string
	Decode(string) (int64, error)
}
