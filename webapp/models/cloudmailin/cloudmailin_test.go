package cloudmailin

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %+v (type %v) - Got %+v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestDecode(t *testing.T) {
	data := `
{
  "headers": {
    "Return-Path": "from@example.com",
    "Received": [
      "by 10.52.90.229 with SMTP id bz5cs75582vdb; Mon, 16 Jan 2012 09:00:07 -0800",
      "by 10.216.131.153 with SMTP id m25mr5479776wei.9.1326733205283; Mon, 16 Jan 2012 09:00:05 -0800",
      "from mail-wi0-f170.google.com (mail-wi0-f170.google.com [209.85.212.170]) by mx.google.com with ESMTPS id u74si9614172weq.62.2012.01.16.09.00.04 (version=TLSv1/SSLv3 cipher=OTHER); Mon, 16 Jan 2012 09:00:04 -0800"
    ],
    "Date": "Mon, 16 Jan 2012 17:00:01 +0000",
    "From": "Message Sender <sender@example.com>",
    "To": "Message Recipient<to@example.co.uk>",
    "Message-ID": "<4F145791.8040802@example.com>",
    "Subject": "Test Subject",
    "Mime-Version": "1.0",
    "Content-Type": "multipart/alternative; boundary=------------090409040602000601080801",
    "Delivered-To": "to@example.com",
    "Received-SPF": "neutral (google.com: 10.0.10.1 is neither permitted nor denied by best guess record for domain of from@example.com) client-ip=10.0.10.1;",
    "Authentication-Results": "mx.google.com; spf=neutral (google.com: 10.0.10.1 is neither permitted nor denied by best guess record for domain of from@example.com) smtp.mail=from@example.com",
    "User-Agent": "Postbox 3.0.2 (Macintosh/20111203)"
  },
  "envelope": {
    "to": "to@example.com",
    "from": "from@example.com",
    "helo_domain": "localhost",
    "remote_ip": "127.0.0.1",
    "recipients": [
      "to@example.com",
      "another@example.com"
    ],
    "spf": {
      "result": "pass",
      "domain": "example.com"
    }
  },
  "plain": "Test with HTML.",
  "html": "<html><head>\n<meta http-equiv=\"content-type\" content=\"text/html; charset=ISO-8859-1\"></head><body\n bgcolor=\"#FFFFFF\" text=\"#000000\">\nTest with <span style=\"font-weight: bold;\">HTML</span>.<br>\n</body>\n</html>",
  "reply_plain": "Message reply if found.",
  "attachments": [
    {
      "content": "dGVzdGZpbGU=",
      "file_name": "file.txt",
      "content_type": "text/plain",
      "size": 8,
      "disposition": "attachment"
    },
    {
      "content": "dGVzdGZpbGU=",
      "file_name": "file.txt",
      "content_type": "text/plain",
      "size": 8,
      "disposition": "attachment"
    }
  ]
}
`
	ret, err := Decode(ioutil.NopCloser(strings.NewReader(data)))
	if err != nil {
		t.Error(err)
	}
	expect(t, ret.Headers.Subject, "Test Subject")
	expect(t, ret.Plain, "Test with HTML.")
	expect(t, ret.Envelope.From, "from@example.com")
}
