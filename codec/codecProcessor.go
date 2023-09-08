package codec

import (
	"encoding/json"
	"sync"
)

type CodecProcessor struct {
}

var (
	instance *CodecProcessor
	once     sync.Once
)

// GetInstance returns a singleton instance of the Factory.
func GetInstance() *CodecProcessor {
	once.Do(func() {
		instance = &CodecProcessor{}

	})
	return instance
}

func (cp *CodecProcessor) Encode(data interface{}) ([]byte, error) {

	encodedData, err := json.Marshal(data)

	return encodedData, err
}

func (cp *CodecProcessor) Decode(encodedData []byte, v any) error {

	err := json.Unmarshal(encodedData, v)
	return err
}
