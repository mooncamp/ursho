package encoding

type Coder interface {
	Encode(int64) (string, error)
	Decode(string) (int64, error)
}
