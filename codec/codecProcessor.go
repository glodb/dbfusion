package codec

import (
	"encoding/json"
	"sync"
)

type codecProcessor struct {
}

var (
	instance *codecProcessor
	once     sync.Once
)

// GetInstance returns a singleton instance of the Factory.
func GetInstance() *codecProcessor {
	once.Do(func() {
		instance = &codecProcessor{}

	})
	return instance
}

func (cp *codecProcessor) Encode(data interface{}) ([]byte, error) {

	encodedData, err := json.Marshal(data)

	return encodedData, err
}

func (cp *codecProcessor) Decode(encodedData []byte, v any) error {

	err := json.Unmarshal(encodedData, v)
	return err
}
