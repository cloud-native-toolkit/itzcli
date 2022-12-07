package prompt

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

type Option struct {
	text      string
	isDefault bool
}

type ValidatorFunc func(*Prompt, string) (bool, error)

type PromptFilterFunc func(prompt *Prompt) bool

type PromptsContext struct {
	answers map[string]string
}

// Set sets the value of the given key
func (c *PromptsContext) Set(key string, value string) {
	c.answers[key] = value
}

func (c *PromptsContext) Lookup(key string) (string, bool) {
	val, ok := c.answers[key]
	return val, ok
}

// NewPromptsContext creates a new PromptsContext
func NewPromptsContext() *PromptsContext {
	return &PromptsContext{
		answers: make(map[string]string),
	}
}

type Prompt struct {
	context       *PromptsContext
	parent        *Prompt
	path          string
	text          string
	options       []Option
	optionHandler ValueGetter
	shortCircuit  PromptFilterFunc
	validator     ValidatorFunc
	subPrompts    []*Prompt
}

func (p *Prompt) GetAnswer(key string) string {
	val, ok := p.context.Lookup(key)
	if ok {
		return val
	}
	return ""
}

func (p *Prompt) LookupAnswer(key string) (string, bool) {
	return p.context.Lookup(key)
}

// VarMap returns a _copy_ of the variable map
func (p *Prompt) VarMap() map[string]string {
	result := make(map[string]string)
	for k, v := range p.context.answers {
		result[k] = v
	}
	return result
}

func (p *Prompt) AddSubPrompt(prompt *Prompt) {
	prompt.parent = p
	prompt.context = p.context
	p.subPrompts = append(p.subPrompts, prompt)
}

func (p *Prompt) String() string {
	return p.text
}

func (p *Prompt) AvailableOptions() []Option {
	allOptions := []Option{}
	if p.optionHandler != nil {
		values, _ := p.optionHandler()
		if len(values) > 0 {
			for _, opt := range values {
				allOptions = append(allOptions, Option{
					text:      opt,
					isDefault: false,
				})
			}
			return append(p.options, allOptions...)
		}
	}
	return p.options
}

func (p *Prompt) Record(answer string) bool {
	p.context.Set(p.path, strings.TrimSpace(answer))
	return true
}

func (p *Prompt) Itr() func() *Prompt {
	// Basically, this is a delegate for closures. Based on different
	// criteria, we want to
	logrus.Trace("Building iterator now...")
	var current *Prompt
	var myIdx int = 0

	return func() *Prompt {
		logrus.Trace("Getting next prompt now...")

		// The first thing we'll do is see if the current one is nil
		// and also make sure we don't have a parent. This is a special
		// case and so set the current to ourselves and return a self-referencing
		// iterator
		if current == nil && p.parent == nil {
			current = p
			return current
		}

		// Now we'll make sure we're valid. If we're not valid, then we
		// also return the self-refernce iterator
		isValid, _ := current.isValid()
		if !isValid {
			return current
		}

		currIdx := 0
		// Now that we're past all this, our goal in life is to return
		// the next child object
		if currIdx < len(current.subPrompts) {
			next := current.subPrompts[currIdx]
			current = next
			currIdx++
			return current
		}
		// but once I'm at the end of my sub items, I need to return my parent's
		// next sub item (my sibling)
		if len(current.subPrompts) == 0 {
			if len(p.subPrompts) > (myIdx + 1) {
				next := p.subPrompts[myIdx+1]
				current = next
				myIdx++
				return current
			}
		}

		return nil
	}
}

func (p *Prompt) isValid() (bool, error) {
	// If it does not have a validator at all, it's just assumed good.
	if p.validator == nil {
		return true, nil
	}
	val, ok := p.context.Lookup(p.path)
	if !ok {
		return false, nil
	}
	return p.validator(p, val)
}

type PromptBuilder struct {
	ctx           *PromptsContext
	path          string
	text          string
	filter        PromptFilterFunc
	options       []Option
	optionFunc    ValueGetter
	validatorFunc ValidatorFunc
}

func (b *PromptBuilder) Context(ctx *PromptsContext) *PromptBuilder {
	b.ctx = ctx
	return b
}

func (b *PromptBuilder) WithLogging() *PromptBuilder {
	logrus.SetLevel(logrus.DebugLevel)
	return b
}

func (b *PromptBuilder) Path(p string) *PromptBuilder {
	b.path = p
	return b
}

func (b *PromptBuilder) Text(t string) *PromptBuilder {
	b.text = t
	return b
}

func (b *PromptBuilder) Textf(format string, a ...interface{}) *PromptBuilder {
	b.text = fmt.Sprintf(format, a...)
	return b
}

func (b *PromptBuilder) AskWhen(filter PromptFilterFunc) *PromptBuilder {
	b.filter = filter
	return b
}

func (b *PromptBuilder) AddOption(text string) *PromptBuilder {
	opt := &Option{
		text:      text,
		isDefault: false,
	}
	b.options = append(b.options, *opt)
	return b
}

func (b *PromptBuilder) AddDefaultOption(text string) *PromptBuilder {
	opt := &Option{
		text:      text,
		isDefault: true,
	}
	b.options = append(b.options, *opt)
	return b
}

func (b *PromptBuilder) WithOptions(optionFunc ValueGetter) *PromptBuilder {
	b.optionFunc = optionFunc
	return b
}

func (b *PromptBuilder) WithValidator(f ValidatorFunc) *PromptBuilder {
	b.validatorFunc = f
	return b
}

func (b *PromptBuilder) Build() (*Prompt, error) {

	if len(b.path) == 0 {
		return nil, fmt.Errorf("there must be a path defined")
	}

	return &Prompt{
		context:       b.ctx,
		path:          b.path,
		text:          b.text,
		shortCircuit:  b.filter,
		options:       b.options,
		validator:     b.validatorFunc,
		optionHandler: b.optionFunc,
	}, nil
}

func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		ctx:     NewPromptsContext(),
		options: []Option{},
	}
}

func Ask(prompt *Prompt, out io.Writer, in io.Reader) error {
	buf := bufio.NewReader(in)
	_, err := out.Write([]byte(fmt.Sprintf("%s  ", prompt.String())))
	if err != nil {
		return err
	}
	answer, _ := buf.ReadString('\n')
	logrus.Tracef("Recording answer: %s", answer)
	prompt.Record(answer)
	return nil
}
