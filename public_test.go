package nanomarkup

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPRequestStruct(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
	}
	enc, err := Marshal(req, nil)
	if err != nil {
		t.Error(err)
	}
	dec := &http.Request{}
	err = Unmarshal(enc, dec, nil)
	if err != nil {
		t.Error(err)
	}
	testStructs(t, req, dec)
}

func TestIndentIndent(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := Marshal(req, nil)
	if err != nil {
		t.Error(err)
		return
	}
	dst := bytes.Buffer{}
	if err = Indent(&dst, enc, "", "    "); err != nil {
		t.Error(err)
		return
	}
	want := `{
    Method GET
    URL {
        Scheme https
        Opaque 
        Host google.com
        Path 
        RawPath 
        OmitHost false
        ForceQuery false
        RawQuery 
        Fragment 
        RawFragment 
    }
    Proto HTTP/1.1
    ProtoMajor 1
    ProtoMinor 1
    Header {
    }
    ContentLength 0
    Close false
    Host google.com
    RemoteAddr 
    RequestURI 
}
`
	out := dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", enc, out, want)
	}
}

func TestIndentPrefix(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := Marshal(req, nil)
	if err != nil {
		t.Error(err)
		return
	}
	dst := bytes.Buffer{}
	if err = Indent(&dst, enc, "##", "  "); err != nil {
		t.Error(err)
		return
	}
	want := `{
##  Method GET
##  URL {
##    Scheme https
##    Opaque 
##    Host google.com
##    Path 
##    RawPath 
##    OmitHost false
##    ForceQuery false
##    RawQuery 
##    Fragment 
##    RawFragment 
##  }
##  Proto HTTP/1.1
##  ProtoMajor 1
##  ProtoMinor 1
##  Header {
##  }
##  ContentLength 0
##  Close false
##  Host google.com
##  RemoteAddr 
##  RequestURI 
##}
`
	out := dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", enc, out, want)
	}
}

func TestIndentMultiline(t *testing.T) {
	// test a string
	s := `testing
a multi
line
value`
	in := "`\n" + s + "\n`\n"
	want := "`\n" + `##testing
##a multi
##line
##value` + "\n##`\n"
	dst := bytes.Buffer{}
	if err := Indent(&dst, []byte(in), "##", "  "); err != nil {
		t.Error(err)
		return
	}
	out := dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", in, out, want)
	}

	// test a struct
	in = "{\nMultiLine `\n" + s + "\n`\n}\n"
	want = "{\n##  MultiLine `\n" + `##testing
##a multi
##line
##value` + "\n##  `\n##}\n"
	dst = bytes.Buffer{}
	if err := Indent(&dst, []byte(in), "##", "  "); err != nil {
		t.Error(err)
		return
	}
	out = dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", in, out, want)
	}
}

func TestIndentComment(t *testing.T) {
	// test a string
	sin := `testing
a multi
line
value`
	swant := `##testing
##a multi
##line
##value`
	comment := "// Test a comment"
	in := fmt.Sprintf("%s\n`\n%s\n`\n", comment, sin)
	want := fmt.Sprintf("%s\n##`\n%s\n##`\n", comment, swant)
	dst := bytes.Buffer{}
	if err := Indent(&dst, []byte(in), "##", "  "); err != nil {
		t.Error(err)
		return
	}
	out := dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", in, out, want)
	}

	// test a struct
	in = fmt.Sprintf("{\n%s\nMultiLine `\n%s\n`\n}\n", comment, sin)
	want = fmt.Sprintf("{\n##  %s\n##  MultiLine `\n%s\n##  `\n##}\n", comment, swant)
	dst = bytes.Buffer{}
	if err := Indent(&dst, []byte(in), "##", "  "); err != nil {
		t.Error(err)
		return
	}
	out = dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", in, out, want)
	}
}

func TestMarshalIndent(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := MarshalIndent(req, "\t", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	want := "{\n" +
		"\t\tMethod GET\n" +
		"\t\tURL {\n" +
		"\t\t\tScheme https\n" +
		"\t\t\tOpaque \n" +
		"\t\t\tHost google.com\n" +
		"\t\t\tPath \n" +
		"\t\t\tRawPath \n" +
		"\t\t\tOmitHost false\n" +
		"\t\t\tForceQuery false\n" +
		"\t\t\tRawQuery \n" +
		"\t\t\tFragment \n" +
		"\t\t\tRawFragment \n" +
		"\t\t}\n" +
		"\t\tProto HTTP/1.1\n" +
		"\t\tProtoMajor 1\n" +
		"\t\tProtoMinor 1\n" +
		"\t\tHeader {\n" +
		"\t\t}\n" +
		"\t\tContentLength 0\n" +
		"\t\tClose false\n" +
		"\t\tHost google.com\n" +
		"\t\tRemoteAddr \n" +
		"\t\tRequestURI \n" +
		"\t}\n"
	if string(enc) != want {
		t.Errorf("[MarshalIndent] in: %s; out: %s; want: %s", enc, string(enc), want)
	}
}

func TestCompact(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := Marshal(req, nil)
	if err != nil {
		t.Error(err)
		return
	}
	ind := bytes.Buffer{}
	if err = Indent(&ind, enc, "\t", " "); err != nil {
		t.Error(err)
		return
	}
	dst := bytes.Buffer{}
	if err = Compact(&dst, ind.Bytes()); err != nil {
		t.Error(err)
		return
	}
	want := `{
Method GET
URL {
Scheme https
Opaque 
Host google.com
Path 
RawPath 
OmitHost false
ForceQuery false
RawQuery 
Fragment 
RawFragment 
}
Proto HTTP/1.1
ProtoMajor 1
ProtoMinor 1
Header {
}
ContentLength 0
Close false
Host google.com
RemoteAddr 
RequestURI 
}
`
	out := dst.String()
	if out != want {
		t.Errorf("[Compact] in: %s; out: %s; want: %s", ind.String(), out, want)
	}
	if out != string(enc) {
		t.Errorf("[Compact] in: %s; out: %s; want: %s", ind.String(), out, string(enc))
	}
}
