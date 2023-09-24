package codec

import (
	"encoding/json"
	"sync"
)

// codecProcessor is a concrete singleton processor class within the codec package. It provides
// support for encoding and decoding data to and from JSON format. This class ensures that only
// one instance of itself is created using the singleton pattern.
type codecProcessor struct {
}

var (
	instance *codecProcessor // Singleton instance of the codecProcessor.
	once     sync.Once       // Once ensures the singleton instance is created only once.
)

// GetInstance is a method of the codecProcessor class that returns a singleton instance
// of the codecProcessor. It uses the sync.Once mechanism to ensure that the instance is
// created only once during the first invocation of this method.
//
// Returns:
//   - *codecProcessor: A pointer to the singleton instance of the codecProcessor.
func GetInstance() *codecProcessor {
	// Use sync.Once to ensure that the instance is created only once.
	once.Do(func() {
		instance = &codecProcessor{}
	})

	// Return the singleton instance.
	return instance
}

// Encode is a method of the codecProcessor class that encodes a given data structure into
// a byte slice using the JSON encoding format.
//
// Parameters:
//   - data: The data structure to be encoded.
//
// Returns:
//   - []byte: The encoded byte slice representing the data in JSON format.
//   - error: An error if encoding fails.
func (cp *codecProcessor) Encode(data interface{}) ([]byte, error) {
	// Use the JSON encoding to encode the given data structure.
	encodedData, err := json.Marshal(data)

	// Return the encoded byte slice and any potential errors.
	return encodedData, err
}

// Decode is a method of the codecProcessor class that decodes a byte slice containing data
// in JSON format back into its original data structure.
//
// Parameters:
//   - encodedData: The byte slice containing the data to be decoded in JSON format.
//   - v: A pointer to the target data structure into which the decoded data will be stored.
//
// Returns:
//   - error: An error if decoding fails.
func (cp *codecProcessor) Decode(encodedData []byte, v any) error {
	// Use JSON decoding to decode the byte slice back into the original data structure.
	err := json.Unmarshal(encodedData, v)

	// Return any potential errors that may occur during decoding.
	return err
}
