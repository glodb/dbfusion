package codecs_test

import (
	"testing"

	"github.com/glodb/dbfusion/codec"
	"github.com/glodb/dbfusion/tests/models"
)

var data map[string]interface{}
var encodedJsonData []byte

func init() {
	data = make(map[string]interface{})
	data["firstName"] = "Aafaq"
	data["email"] = "aafaqzahid9@gmail.com"
	data["userName"] = "aafaqzahid9"
	data["password"] = "change-me"
}

func TestJsonEncode(t *testing.T) {
	var err error
	encodedJsonData, err = codec.GetInstance().Encode(data)
	if err != nil {
		t.Errorf("Error in encoding json %v", err)
	}
}

func TestJsonDecode(t *testing.T) {
	newData := models.UserTest{}
	err := codec.GetInstance().Decode(encodedJsonData, &newData)
	if err != nil {
		t.Errorf("Error in decoding json %v", err)
	}
}
