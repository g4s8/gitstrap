package spec

import (
	"fmt"
	"strings"
)

type Readme struct {
	Selector struct {
		Repository string `yaml:"repository"`
	} `yaml:"selector"`
	Title    string `yaml:"title,omitempty"`
	Abstract string `yaml:"abstract,omitempty"`
	Topics   []struct {
		Heading string `yaml:"heading"`
		Body    string `yaml:"body"`
	} `yaml:"topics,omitempty"`
}

func (s *Readme) String() string {
	sb := new(strings.Builder)
	if s.Title != "" {
		sb.WriteString(fmt.Sprintf("# %s\n\n", s.Title))
	}
	if s.Abstract != "" {
		sb.WriteString(s.Abstract)
		sb.WriteString("\n\n")
	}
	for _, topic := range s.Topics {
		sb.WriteString(fmt.Sprintf("## %s\n\n", topic.Heading))
		sb.WriteString(topic.Body)
		sb.WriteString("\n\n")
	}
	return sb.String()
}
