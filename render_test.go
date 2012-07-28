package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestRenderMarkdown(t *testing.T) {
	m := map[string]string{
		"Hello.":        "<p>Hello.</p>",
		"Hi **there**!": "<p>Hi <strong>there</strong>!</p>",
	}

	for input, expected := range m {
		got := strings.TrimSpace(string(RenderMarkdown([]byte(input))))
		if got != expected {
			t.Errorf("'%s': expected '%s', got '%s'", input, expected, got)
		} else {
			t.Logf("'%s': '%s' (OK)", input, got)
		}
	}
}

func TestTemplateComposition(t *testing.T) {
	m := map[string]string{
		"a.template": "hello {{user}}",
		"b.template": "<p>[[a.template]]</p>",
		"c.template": "<body>[[b.template]]</body>",
	}

	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("using tempDir %s", tempDir)
	defer os.RemoveAll(tempDir)

	for filename, body := range m {
		f, err := os.Create(tempDir + "/" + filename)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.WriteString(body); err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	ctx := Context{"user": "friend"}
	buf, err := RenderTemplate(tempDir, "c.template", ctx)
	if err != nil {
		t.Fatal(err)
	}

	got := string(buf)
	expected := "<body><p>hello friend</p></body>"
	if got != expected {
		t.Fatalf("expected '%s', got '%s'", expected, got)
	} else {
		t.Logf("%s (OK)", got)
	}
}
