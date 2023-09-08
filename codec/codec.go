package codec

type Codec interface {
	Encode(data interface{}) ([]byte, error)
	Decode([]byte, any) error
}
