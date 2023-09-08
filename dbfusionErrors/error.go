package dbfusionErrors

import "errors"

var (
	ErrDBTypeNotSupported       = errors.New("This DB Type is not supported")
	ErrConnectionNotAvailable   = errors.New("Connection which you are trying to close is not available")
	ErrUriRequiredForConnection = errors.New("Uri is required for the connection")
	ErrCacheNotConnected        = errors.New("Cache object is passed but it is not connected")
	ErrStructOrMapRequired      = errors.New("Struct or map[struct]interface is required for addition into dbs")
	ErrCacheIndexesIncreased    = errors.New("10 Cache indexes are allowed in a single struct")
	ErrCacheUniqueKeysIncreased = errors.New("A composite key of 5 variables is allowed at this moment")
	ErrCodecFormatNotSupported  = errors.New("Current selected format is not supported")
)
