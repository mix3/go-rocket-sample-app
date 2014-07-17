package mailgun

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	os.Setenv("MAILGUN_APIKEY", "fake API key")
	os.Setenv("MAILGUN_DOMAIN", "testdomain.org")
	test_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatal("Error when parsing the form: ", err)
		}

		recvToExpected := map[string]string{
			r.FormValue("from"):    "name <from@example.com>",
			r.FormValue("to"):      "to@example.com",
			r.FormValue("subject"): "subject",
			r.FormValue("text"):    "body",
		}
		for received, expected := range recvToExpected {
			assert.Equal(t, expected, received)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Queued. Thank you.", "id":"fake_id_here@mailgun.org"}`)
	}))
	defer test_server.Close()
	client := NewClient()
	client.Hostname = test_server.URL
	res, err := client.Send("name", "from@example.com", "to@example.com", "subject", "body")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `{"message": "Queued. Thank you.", "id":"fake_id_here@mailgun.org"}`, res)
}
