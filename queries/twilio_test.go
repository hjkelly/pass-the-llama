package queries

import (
	"testing"
)

func TestSendSms(t *testing.T) {
	err := SendSms("+15005550006", "testing Pass the Llama's SendSms function")
	if err != nil {
		t.Error(err.Error())
	}
}
