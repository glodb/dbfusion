package codec

// Codec is an interface defined in the codecw package for encoding and decoding data structures
// to facilitate storage and retrieval in a cache. Implementations of this interface provide methods
// for encoding data into byte slices and decoding byte slices back into their original data form.
type Codec interface {
	// Encode encodes a given data structure into a byte slice.
	//
	// Parameters:
	//   - data: The data structure to be encoded.
	//
	// Returns:
	//   - []byte: The encoded byte slice representing the data.
	//   - error: An error if encoding fails.
	Encode(data interface{}) ([]byte, error)

	// Decode decodes a byte slice back into its original data structure.
	//
	// Parameters:
	//   - encodedData: The byte slice to be decoded.
	//   - target: A pointer to the target data structure into which the decoded data will be stored.
	//
	// Returns:
	//   - error: An error if decoding fails.
	Decode(encodedData []byte, target interface{}) error
}
