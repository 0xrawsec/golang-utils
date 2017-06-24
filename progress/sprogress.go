package progress

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
)

const ellipse string = "[...]"

// ErrBadSize is generated if the string size of the progress is not valid
var ErrBadSize error = errors.New("ErrBadSize: progress message size not allowed")

// Progress structure
type Progress struct {
	pre     string
	state   string
	message string
	sOutput int
}

// New Progress structure
func New(optionals ...int) Progress {
	if len(optionals) > 0 {
		if optionals[0] < 0 {
			panic(ErrBadSize)
		}
		return Progress{pre: "Progress", state: "|", sOutput: optionals[0]}
	}
	return Progress{pre: "Progress", state: "|"}
}

// SetPre set the prologue part of the progress message
func (p *Progress) SetPre(newPre string) {
	p.pre = newPre
}

//Update the status of the Progress
func (p *Progress) Update(message string) {
	switch p.state {
	case "|":
		p.state = "/"
	case "/":
		p.state = "-"
	case "-":
		p.state = "\\"
	case "\\":
		p.state = "|"
	}
	if message != "" {
		p.message = message
	}
}

func (p *Progress) String() string {
	strProgressInvariant := fmt.Sprintf("%s: %s", p.pre, p.state)
	if len(strProgressInvariant)+len(p.message)+1 <= p.sOutput {
		return fmt.Sprintf("%s %s", strProgressInvariant, p.message)
	}
	sNewMessage := p.sOutput - (len(strProgressInvariant) + len(ellipse) + 1)
	if sNewMessage < 0 {
		return strProgressInvariant
	}
	limit := int(math.Floor(float64(sNewMessage) / 2))
	p.message = p.message[:limit] + ellipse + p.message[len(p.message)-limit:]
	if len(fmt.Sprintf("%s %s", strProgressInvariant, p.message)) > p.sOutput {
		panic(ErrBadSize)
	}
	return fmt.Sprintf("%s %s", strProgressInvariant, p.message)
}

// Print Progress structure on stderr or stdout if the first optional argument
// is true
func (p *Progress) Print(optionals ...bool) {
	stream := os.Stderr
	if len(optionals) > 0 {
		if optionals[0] {
			stream = os.Stdout
		}
	}
	f := bufio.NewWriter(stream)
	defer f.Flush()
	f.WriteString(fmt.Sprintf("% *s\r", p.sOutput, ""))
	f.WriteString(p.String() + "\r")
}
