package email

import "testing"

func TestSendWithAttachment(t *testing.T) {
	err := SendWithAttachment(DefaultFilter,
		"TEST", "TEST-TEST", "Nothing Here<br/><br/>test<br/><br/>mr",
		[]string{"yuan.ren@ifchange.com"}, []string{"yuan.ren@ifchange.com"},
		map[string][]byte{
			"a.html": []byte("<html></html>"),
		})
	if err != nil {
		t.Fatal(err)
	}
}
