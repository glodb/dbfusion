package codecs_test

import (
	"testing"

	"github.com/glodb/dbfusion/codec"
	"github.com/glodb/dbfusion/tests/models"
)

var data map[string]interface{}
var encodedJsonData []byte

// init is called before the test functions and initializes the test data.
func init() {
	data = make(map[string]interface{})
	data["firstName"] = "Aafaq"
	data["email"] = "aafaqzahid9@gmail.com"
	data["userName"] = "aafaqzahid9"
	data["password"] = "change-me"
}

// TestJsonEncode tests the JSON encoding functionality.
func TestJsonEncode(t *testing.T) {
	var err error

	// Encode the test data into JSON.
	encodedJsonData, err = codec.GetInstance().Encode(data)

	// Check if there was an error during encoding and report it.
	if err != nil {
		t.Errorf("Error in encoding JSON: %v", err)
	}
}

// TestJsonDecode tests the JSON decoding functionality.
func TestJsonDecode(t *testing.T) {
	// Create a new instance of the UserTest struct to decode into.
	newData := models.UserTest{}

	// Decode the JSON data into the newData variable.
	err := codec.GetInstance().Decode(encodedJsonData, &newData)

	// Check if there was an error during decoding and report it.
	if err != nil {
		t.Errorf("Error in decoding JSON: %v", err)
	}
}
