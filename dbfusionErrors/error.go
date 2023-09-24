package dbfusionErrors

import "errors"

var (
	ErrDBTypeNotSupported       = errors.New("This DB Type is not supported")
	ErrConnectionNotAvailable   = errors.New("Connection which you are trying to close is not available")
	ErrUriRequiredForConnection = errors.New("Uri is required for the connection")
	ErrCacheNotConnected        = errors.New("Cache object is passed but it is not connected")
	ErrStringMapRequired        = errors.New("Struct or map[string]interface for db operations")
	ErrCacheIndexesIncreased    = errors.New("10 Cache indexes are allowed in a single struct")
	ErrCacheUniqueKeysIncreased = errors.New("A composite key of 5 variables is allowed at this moment")
	ErrCodecFormatNotSupported  = errors.New("Current selected format is not supported")
	ErrEntityNameRequired       = errors.New("In case of map[string]interface entity name is required")
	ErrNoValidCacheFound        = errors.New("No valid cache is found to process this hook")
	ErrSQLQueryTypeNotSupported = errors.New("Map or built in types are supported for query")
	ErrSQLQueryNoRecordFound    = errors.New("No Record found in sql query")
	ErrInvalidType              = errors.New("This type is not supported by the library yet")
	ErrNoRecordFound            = errors.New("No record found for the query")
)
