package dbfusionErrors

import (
	"errors"
)

// ErrDBTypeNotSupported is returned when the provided database type is not supported by the library.
var ErrDBTypeNotSupported = errors.New("This DB Type is not supported")

// ErrConnectionNotAvailable is returned when attempting to close a connection that is not available.
var ErrConnectionNotAvailable = errors.New("Connection which you are trying to close is not available")

// ErrUriRequiredForConnection is returned when a URI is required for establishing a connection, but it is missing.
var ErrUriRequiredForConnection = errors.New("Uri is required for the connection")

// ErrCacheNotConnected is returned when a cache object is passed, but it is not connected.
var ErrCacheNotConnected = errors.New("Cache object is passed but it is not connected")

// ErrStringMapRequired is returned when a struct or map[string]interface{} is required for database operations.
var ErrStringMapRequired = errors.New("Struct or map[string]interface for db operations")

// ErrCacheIndexesIncreased is returned when attempting to use more than 10 cache indexes in a single struct.
var ErrCacheIndexesIncreased = errors.New("10 Cache indexes are allowed in a single struct")

// ErrCacheUniqueKeysIncreased is returned when attempting to create a composite key with more than 5 variables.
var ErrCacheUniqueKeysIncreased = errors.New("A composite key of 5 variables is allowed at this moment")

// ErrCodecFormatNotSupported is returned when the selected encoding format is not supported by the library.
var ErrCodecFormatNotSupported = errors.New("Current selected format is not supported")

// ErrEntityNameRequired is returned when an entity name is required, especially in the case of map[string]interface{}.
var ErrEntityNameRequired = errors.New("In case of map[string]interface entity name is required")

// ErrNoValidCacheFound is returned when no valid cache is found to process a cache hook.
var ErrNoValidCacheFound = errors.New("No valid cache is found to process this hook")

// ErrSQLQueryTypeNotSupported is returned when the query type is not supported. Supported types include map and built-in types.
var ErrSQLQueryTypeNotSupported = errors.New("Map or built-in types are supported for query")

// ErrSQLQueryNoRecordFound is returned when no records are found for an SQL query.
var ErrSQLQueryNoRecordFound = errors.New("No Record found in sql query")

// ErrInvalidType is returned when an unsupported type is encountered.
var ErrInvalidType = errors.New("This type is not supported by the library yet")

// ErrNoRecordFound is returned when no records are found for a query.
var ErrNoRecordFound = errors.New("No record found for the query")
