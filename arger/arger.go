// Package arger
// Default value ::
//   -- stdin
//   ** stdout
//   @@ stderr
//   !! needed

package arger

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func convInt(str string) (int, error) {
	l := len(str)
	if l == 0 {
		return 0, fmt.Errorf("Empty string input%s", "!!!")
	}

	multi := 1
	switch str[l-1] {
	case 'K', 'k':
		multi = 1024
		str = str[:l-1]
	case 'M', 'm':
		multi = 1024 * 1024
		str = str[:l-1]
	case 'G', 'g':
		multi = 1024 * 1024 * 1024
		str = str[:l-1]
	case 'T', 't':
		multi = 1024 * 1024 * 1024 * 1024
		str = str[:l-1]
	case 'P', 'p':
		multi = 1024 * 1024 * 1024 * 1024 * 1024
		str = str[:l-1]
	}
	if len(str) == 0 {
		return multi, nil
	}

	n, err := strconv.Atoi(str)
	if err != nil {
		return n, err
	}
	return n * multi, nil
}

type value interface {
	String() string
	Type() string
	Get() interface{}
	Set(string) error
}

type boolValue bool

func (b *boolValue) String() string   { return fmt.Sprintf("%v", *b) }
func (b *boolValue) Type() string     { return "bool" }
func (b *boolValue) Get() interface{} { return bool(*b) }
func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	*b = boolValue(v)
	return err
}

type intValue int

func (i *intValue) String() string   { return fmt.Sprintf("%v", *i) }
func (i *intValue) Type() string     { return "int" }
func (i *intValue) Get() interface{} { return int(*i) }
func (i *intValue) Set(s string) error {
	v, err := convInt(s)
	*i = intValue(v)
	return err
}

type stringValue string

func (s *stringValue) String() string   { return fmt.Sprintf("%v", *s) }
func (s *stringValue) Type() string     { return "string" }
func (s *stringValue) Get() interface{} { return string(*s) }
func (s *stringValue) Set(v string) error {
	*s = stringValue(v)
	return nil
}

type floatValue float64

func (f *floatValue) String() string   { return fmt.Sprintf("%v", *f) }
func (f *floatValue) Type() string     { return "float" }
func (f *floatValue) Get() interface{} { return float64(*f) }
func (f *floatValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	*f = floatValue(v)
	return err
}

type parameter struct {
	Name   string
	Desc   string
	sterm  string
	lterm  string
	value  value
	defval value
}

func (p *parameter) Type() string          { return p.defval.Type() }
func (p *parameter) String() string        { return p.defval.String() }
func (p *parameter) DefValue() interface{} { return p.defval.Get() }
func (p *parameter) Get() interface{}      { return p.value.Get() }
func (p *parameter) Set(s string) error    { return p.value.Set(s) }

func newParameter(name, sterm, lterm, desc string, defval interface{}) *parameter {
	var val value
	switch defval.(type) {
	case int:
		v := defval.(int)
		d := new(intValue)
		*d = intValue(v)
		val = d
	case bool:
		v := defval.(bool)
		d := new(boolValue)
		*d = boolValue(v)
		val = d
	case string:
		v := defval.(string)
		d := new(stringValue)
		*d = stringValue(v)
		val = d
	case float64:
		v := defval.(float64)
		d := new(floatValue)
		*d = floatValue(v)
		val = d
	default: // unsupported type
		panic(fmt.Errorf("Unsupport Type: %T(%v)", defval, defval))
	}

	defVal := val
	return &parameter{
		Name:   name,
		Desc:   desc,
		sterm:  sterm,
		lterm:  lterm,
		value:  val,
		defval: defVal,
	}
}

// Parser (Name, SubName, usage, Args)
type Parser struct {
	Name    string   // Command Name
	SubName string   // subcommand Name
	usage   string   // show custom usage
	Args    []string // record args
	notes   []string // record note info
	paras   map[string]*parameter
}

// SetUsage set custom usage
func (a *Parser) SetUsage(s string) {
	a.usage = s
}

// Note add note in ArgPareser
func (a *Parser) Note(s string) {
	a.notes = append(a.notes, s)
}

// Get get para name value
func (a *Parser) Get(name string) interface{} {
	_, ok := a.paras[name]
	if !ok {
		return nil
	}
	return a.paras[name].Get()
}

func (a *Parser) getParameter(term string) *parameter {
	if len(a.paras) == 0 {
		return nil
	}

	for _, p := range a.paras {
		if p.sterm == term || p.lterm == term {
			return p
		}
	}
	return nil
}

// Add add a parameter
func (a *Parser) Add(name, sterm, lterm, describe string, defval interface{}) {
	a.paras[name] = newParameter(name, sterm, lterm, describe, defval)
}

func getParaLength(paras map[string]*parameter) int {
	length := 0
	for _, p := range paras {
		l := 0
		if p.sterm != "" {
			l += len(p.sterm)
		}
		if p.lterm != "" {
			l += len(p.lterm)
		}
		if p.sterm != "" && p.lterm != "" {
			l += 4
		}
		if l >= length {
			length = l + 1
		}
	}
	return length
}

func getOffsetLength(p *parameter) int {
	l := 0
	if p.sterm != "" {
		l += len(p.sterm)
	}
	if p.lterm != "" {
		l += len(p.lterm)
	}
	if p.sterm != "" && p.lterm != "" {
		l += 4
	}
	return l
}

// Usage show Parser usage
func (a *Parser) Usage() {
	hasOpts := false
	if len(a.paras) > 0 {
		hasOpts = true
	}
	optstr := ""
	if hasOpts {
		optstr = " [opts]"
	}

	if a.usage != "" {
		fmt.Fprintln(os.Stderr, a.usage)
	} else if a.SubName == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s%s filename ...\n", a.Name, optstr)
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s %s%s filename ...\n", a.Name, a.SubName, optstr)
	}

	if len(a.paras) == 0 { // not opts parameters
		return
	}

	if hasOpts {
		length := getParaLength(a.paras)
		fmt.Fprintln(os.Stderr, "  Options:")

		// 	sort opt in order not random
		keys := []string{}
		for key := range a.paras {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// output options
		for _, key := range keys {
			p := a.paras[key]
			offset := getOffsetLength(p)
			dtype := p.Type()
			typestr := ""
			defstr := ""
			if dtype == "bool" {
				typestr = strings.Repeat(" ", length-offset+7)
			} else {
				typestr = strings.Repeat(" ", length-offset) + dtype + strings.Repeat(" ", 7-len(dtype))
				switch {
				case p.Type() == "string" && p.String() == "!!":
					defstr = "(Needed!)"
				case p.Type() == "string" && p.String() == "--":
					defstr = "(Default: Stdin)"
				case p.Type() == "string" && p.String() == "**":
					defstr = "(Default: Stdout)"
				case p.Type() == "string" && p.String() == "@@":
					defstr = "(Default: Stderr)"
				case p.Type() == "string" && p.String() == "":
					defstr = ""
				case p.Type() == "int" && p.String() == "0":
					defstr = ""
				case p.Type() == "float" && p.String() == "0.0":
					defstr = ""
				default:
					defstr = "(Default:" + p.String() + ")"
				}
			}
			if p.lterm != "" && p.sterm != "" {
				fmt.Fprintf(os.Stderr, "\t%s or %s %s %s %s\n", p.sterm, p.lterm, typestr, p.Desc, defstr)
			} else if p.sterm != "" {
				fmt.Fprintf(os.Stderr, "\t%s %s %s %s\n", p.sterm, typestr, p.Desc, defstr)
			} else if p.lterm != "" {
				fmt.Fprintf(os.Stderr, "\t%s %s %s %s\n", p.lterm, typestr, p.Desc, defstr)
			}
		}
	}
	// output notes
	if len(a.notes) > 0 {
		fmt.Fprintln(os.Stderr, "\nNotes::")
		for i, l := 0, len(a.notes); i < l; i++ {
			fmt.Fprintln(os.Stderr, "   ", a.notes[i])
		}
	}
}

// Parse parse arguments
func (a *Parser) Parse(args ...string) error {
	for i, l := 0, len(args); i < l; i++ {
		arg := args[i]
		if arg[0] == '-' && len(arg) > 1 { // treat "-" as stdin
			p := a.getParameter(arg)
			if p == nil {
				return fmt.Errorf("No such option: %s", arg)
			}
			if p.Type() == "bool" {
				p.Set("true")
			} else {
				i++
				if i >= l {
					return fmt.Errorf("Option: %s need a value (type of %s), but nothing given", arg, p.Type())
				}
				err := p.Set(args[i])
				if err != nil {
					return fmt.Errorf("Option: %s need type of %s value, but given %s", arg, p.Type(), args[i])
				}
			}
		} else {
			a.Args = append(a.Args, arg)
		}
	}
	return nil
}

// New New(Command, [SubCommand])
func New(args ...string) *Parser {
	name := "Unkown"
	subName := ""
	l := len(args)
	switch {
	case l >= 2:
		subName = args[1]
		fallthrough
	case l == 1:
		name = args[0]
	}
	return &Parser{
		Name:    name,
		SubName: subName,
		Args:    []string{},
		paras:   make(map[string]*parameter),
	}
}
