package gen

import (
	"bufio"
	"bytes"
	_ "embed"
	"go/format"
	"io"
	"regexp"
	"strings"
	"text/template"
)

var matchGoPackage = regexp.MustCompile("option go_package = \"(?P<package>.*);(?P<name>.*)\";")
var matchMessage = regexp.MustCompile("message (?P<message>.*) {")

type generator struct {
	Topic   string
	Name    string
	Message string
}

type protoData struct {
	Pkg  string
	Name string

	Generators []generator
}

func RenderReader(r io.Reader) (string, bool, error) {
	data, err := parseReader(r)
	if err != nil {
		return "", false, err
	}

	if len(data.Generators) == 0 {
		return "", false, nil
	}

	rendered, err := render(data)
	if err != nil {
		return "", false, nil
	}
	return rendered, true, nil
}

func parseReader(r io.Reader) (protoData, error) {
	var currentGen generator
	searchMessage := false
	var data protoData

	scan := bufio.NewReader(r)

	for {
		l, _, err := scan.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return protoData{}, err
		}
		line := string(l)
		if sb := matchGoPackage.FindStringSubmatch(line); len(sb) > 0 {
			pkg := sb[matchGoPackage.SubexpIndex("package")]
			name := sb[matchGoPackage.SubexpIndex("name")]
			data.Name = name
			data.Pkg = pkg
			break
		}
	}

	for {
		l, _, err := scan.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return protoData{}, err
		}

		line := string(l)
		if searchMessage {
			if mb := matchMessage.FindStringSubmatch(line); len(mb) > 0 {
				currentGen.Message = mb[matchMessage.SubexpIndex("message")]
				searchMessage = false
				data.Generators = append(data.Generators, currentGen)
				continue
			}
		}

		if strings.HasPrefix(line, "// pubgen") {
			name, topic, ok := parsePubgenLine(line)
			if ok {
				currentGen.Name = name
				currentGen.Topic = topic
				searchMessage = true
			}
			continue
		}
	}

	return data, nil
}

//go:embed gettmp.txt
var tmp string

func render(data protoData) (string, error) {
	tmpl, err := template.New("test").Parse(tmp)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	formattedOut, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}

	return string(formattedOut), nil
}

var matchName = regexp.MustCompile("name:([^ ]*)")
var matchTopic = regexp.MustCompile("topic:([^ ]*)")

func parsePubgenLine(line string) (name string, topic string, ok bool) {
	nb := matchName.FindStringSubmatch(line)
	if len(nb) == 0 {
		ok = false
		return
	}
	tb := matchTopic.FindStringSubmatch(line)
	if len(tb) == 0 {
		ok = false
		return
	}

	name = nb[1]
	topic = tb[1]
	ok = true

	return
}
