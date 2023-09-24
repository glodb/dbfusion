// Package codec provides encoding and decoding functionality for setting values in a cache.
// Currently, it supports JSON encoding. Additional encoding methods may be added in the future.
//
// The primary purpose of this package is to encode Go data structures into byte slices for storage
// in a cache and decode them back into their original form when retrieved from the cache.
//
// Supported Encoding Formats:
//   - JSON: This package currently supports encoding and decoding data in JSON format.
//
// Usage:
// To use this package, import it into your Go code and utilize the provided functions for encoding
// and decoding data before storing or retrieving it from a cache.
//
// Example:
//   // Import the codec package
//   import "github.com/globdb/dbfusion/codec"
//
//   // Encode a Go data structure into a byte slice
//   encodedData, err := codec.GetInstance().Encode(data)
//   if err != nil {
//       // Handle the error
//   }
//
//   // Decode a byte slice from the cache back into the original data structure
//   err := codec.GetInstance().Decode(encodedData, &data)
//   if err != nil {
//       // Handle the error
//   }
//
// Note:
// This package is designed to be extensible, and more encoding methods may be added in the future
// to support various data storage and retrieval needs.
package codec
