package parser

import (
	"bytes"
	"fmt"
	"github.com/agaffney/crapsh/parser/ast"
	parser_input "github.com/agaffney/crapsh/parser/input"
	"github.com/agaffney/crapsh/parser/lexer"
	"github.com/agaffney/crapsh/parser/rules/grammar"
	"github.com/agaffney/crapsh/util"
	"io"
)

type Parser struct {
	input       parser_input.Input
	stack       *Stack
	buf         *bytes.Buffer
	tokenBuf    []*lexer.Token
	tokenIdx    int
	commandChan chan ast.Node
	errorChan   chan error
	lexer       *lexer.Lexer
}

func NewParser() *Parser {
	parser := &Parser{}
	parser.lexer = lexer.New()
	return parser
}

func (p *Parser) Reset() {
	p.commandChan = make(chan ast.Node)
	p.errorChan = make(chan error)
	p.stack = &Stack{parser: p}
	p.buf = bytes.NewBuffer(nil)
	p.lexer.Reset()
}

func (p *Parser) Parse(input parser_input.Input) {
	p.Reset()
	p.lexer.Start(input)
	/*
		for {
			token, err := p.lexer.ReadToken()
			fmt.Printf("token = %#v\n", token)
			if err != nil {
				if err == io.EOF {
					if input.IsAvailable() {
						// Restart the lexer to keep pulling from the input
						p.lexer.Start(input)
					} else {
						os.Exit(0)
					}
				} else {
					log.Fatal(err)
				}
			}
		}
		return
	*/
	p.tokenBuf = make([]*lexer.Token, 0)
	p.tokenIdx = -1
	go func() {
		for {
			fmt.Printf("p.stack = %#v\n", p.stack)
			// Reset the hint stack
			p.stack.Reset()
			// Start parsing with "root" element
			ok, err := p.parseRule("BasicCommand", true)
			if err != nil {
				p.errorChan <- err
				close(p.errorChan)
				break
			}
			// EOF
			if !ok {
				break
			}
		}
		close(p.commandChan)
	}()
}

/*
func (p *Parser) GetError() error {
	select {
	case err := <-p.errorChan:
		return err
	default:
		return nil
	}
}
*/

func (p *Parser) GetCommand() ast.Node {
	cmd := <-p.commandChan
	if cmd != nil {
		return cmd
	} else {
		return nil
	}
}

// Handles an individual parser hint
func (p *Parser) parseHandleHint(hint *grammar.ParserHint) (bool, error) {
	util.DumpObject(hint, "parseHandleHint(): hint = ")
	var err error
	var count int
	var ok bool
	var origTokenIdx int
	for {
		ok = false
		origTokenIdx = p.getTokenIdx()
		switch {
		case hint.Type == grammar.HINT_TYPE_RULE:
			ok, err = p.parseRule(hint.Name, false)
		case hint.Type == grammar.HINT_TYPE_ANY:
			ok, err = p.parseAny(hint.Members)
		case hint.Type == grammar.HINT_TYPE_GROUP:
			ok, err = p.parseGroup(hint.Members, false)
		case hint.Type == grammar.HINT_TYPE_TOKEN:
			ok, err = p.parseToken(hint)
		default:
			return ok, fmt.Errorf("Unhandled hint type: %d\n", hint.Type)
		}
		util.DumpObject(ok, "parseHandleHint(): ok = ")
		if err != nil {
			return ok, err
		}
		if !ok {
			fmt.Printf("parseHandleHint(): !ok, hint.Many=%t, count=%d\n", hint.Many, count)
			p.setTokenIdx(origTokenIdx)
			if hint.Optional || (hint.Many && count > 0) {
				return true, nil
			}
			return false, nil
		}
		if !hint.Many {
			break
		}
		count++
	}
	return true, nil
}

func (p *Parser) parseToken(hint *grammar.ParserHint) (bool, error) {
	util.DumpObject(hint, "parseToken(): hint = ")
	token, err := p.nextToken()
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	util.DumpObject(p.curToken(), "parseToken(): curToken = ")
	tokenMatch := false
	for _, hint_token := range hint.TokenTypes {
		if hint_token == token.Type {
			tokenMatch = true
			break
		}
	}
	if tokenMatch {
		p.stack.Cur().astNode.AddToken(token)
		return true, nil
	}
	return false, nil
}

// Handles a 'group' parser hint
// Succeeds if all parser hints match
func (p *Parser) parseGroup(hints []*grammar.ParserHint, updateHintIdx bool) (bool, error) {
	for idx, hint := range hints {
		util.DumpObject(hint, "parseGroup(): hint = ")
		if updateHintIdx {
			p.stack.Cur().hintIdx = idx
		}
		ok, err := p.parseHandleHint(hint)
		if err != nil {
			return ok, err
		}
		if !ok {
			return ok, nil
		}
	}
	return true, nil
}

// Handles a 'rule' parser hint
func (p *Parser) parseRule(ruleName string, send_channel bool) (bool, error) {
	//var parentEndToken *grammar.ParserHint
	rule := grammar.GetRule(ruleName)
	if rule == nil {
		return false, nil
	}
	util.DumpObject(rule, "parseRule(): rule = ")
	p.stack.Add(rule)
	e := ast.NewNode()
	//fmt.Printf("%#v\n", e)
	/*
		if p.stack.Cur().rule.AstFunc != nil {
			foo := p.stack.Cur().rule.AstFunc(e)
			//fmt.Printf("%s\n", foo)
			p.stack.Cur().astNode = foo
		}
	*/
	//p.stack.Cur().parentEndToken = parentEndToken
	p.stack.Cur().astNode = e
	ok, err := p.parseGroup(rule.ParserHints, true)
	if err != nil {
		return false, err
	}
	p.stack.Remove()
	if ok {
		if p.stack.Cur() != nil {
			p.stack.Cur().astNode.AddChild(e)
		} else if send_channel {
			util.DumpJson(e, "sending root element:\n")
			p.commandChan <- e
		}
	}
	return ok, nil
}

// Handles an 'any' parser hint
// Succeeds if any parser hints match
func (p *Parser) parseAny(hints []*grammar.ParserHint) (bool, error) {
	for _, hint := range hints {
		util.DumpObject(hint, "parseAny(): hint = ")
		ok, err := p.parseHandleHint(hint)
		if err != nil {
			return ok, err
		}
		if ok {
			return ok, nil
		}
	}
	return false, nil
}

func (p *Parser) getTokenIdx() int {
	return p.tokenIdx
}

func (p *Parser) setTokenIdx(idx int) {
	p.tokenIdx = idx
	if idx >= 0 {
		fmt.Printf("setTokenIdx(%d): curToken = %#v\n", idx, p.curToken())
	}
}

func (p *Parser) curToken() *lexer.Token {
	return p.tokenBuf[p.tokenIdx]
}

func (p *Parser) prevToken() {
	if p.tokenIdx > 0 {
		p.tokenIdx--
	}
}

func (p *Parser) nextToken() (*lexer.Token, error) {
	if p.tokenIdx < len(p.tokenBuf)-1 {
		p.tokenIdx++
		return p.curToken(), nil
	} else {
		token, _ := p.lexer.ReadToken()
		/*
			if err != nil {
				return nil, err
			}
		*/
		if token == nil {
			return nil, io.EOF
		}
		p.tokenBuf = append(p.tokenBuf, token)
		p.tokenIdx++
		return token, nil
	}
}
