package parser

import (
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
	p.lexer.Reset()
}

func (p *Parser) Start(input parser_input.Input) {
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
			//fmt.Printf("p.stack = %#v\n", p.stack)
			// Reset the hint stack
			p.stack.Reset()
			// Start parsing with "root" element
			ok, err := p.parseRule(`complete_command`, true) // "BasicCommand"
			if err != nil {
				p.errorChan <- err
				close(p.errorChan)
				break
			}
			// EOF (?)
			if !ok {
				break
			}
		}
		close(p.commandChan)
	}()
}

func (p *Parser) GetCommand() (ast.Node, error) {
	select {
	case err := <-p.errorChan:
		return nil, err
	case cmd := <-p.commandChan:
		return cmd, nil
	}
}

// Handles an individual parser hint
func (p *Parser) parseHandleHint(hint *grammar.ParserHint) (bool, error) {
	//util.DumpObject(hint, "parseHandleHint(): hint = ")
	var err error
	var count int
	var ok bool
	var origTokenIdx int
	for {
		ok = false
		origTokenIdx = p.getTokenIdx()
		switch {
		case hint.Type == grammar.HINT_TYPE_RULE:
			ok, err = p.parseRule(hint.RuleName, false)
		case hint.Type == grammar.HINT_TYPE_ANY:
			ok, err = p.parseAny(hint.Members)
		case hint.Type == grammar.HINT_TYPE_GROUP:
			ok, err = p.parseGroup(hint.Members, false)
		case hint.Type == grammar.HINT_TYPE_TOKEN:
			ok, err = p.parseToken(hint)
		default:
			return ok, fmt.Errorf("Unhandled hint type: %d\n", hint.Type)
		}
		//util.DumpObject(ok, "parseHandleHint(): ok = ")
		if err != nil {
			//fmt.Printf("parseHandleHint(): err = %s\n", err.Error())
			return ok, err
		}
		if !ok {
			//fmt.Printf("parseHandleHint(): !ok, hint.Many=%t, hint.Optional=%t, count=%d\n", hint.Many, hint.Optional, count)
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
	//util.DumpObject(hint, "parseToken(): hint = ")
	token, err := p.nextToken()
	if err != nil {
		if err == io.EOF {
			return false, nil
		}
		return false, err
	}
	//util.DumpObject(p.curToken(), "parseToken(): curToken = ")
	tokenType := p.classifyToken(token, hint)
	if reservedRule := p.checkTokenIsReserved(tokenType); reservedRule != nil {
		if !reservedRule.DisallowReservedFollow {
			p.stack.Cur().allowNextWordReserved = true
		}
	}
	tokenMatch := false
	for _, hint_token := range hint.TokenTypes {
		if hint_token == tokenType {
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
		//util.DumpObject(hint, fmt.Sprintf("parseGroup(): hint[%d] = ", idx))
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
		// Set the "final" flag on the current stack entry if we've matched any hints
		if updateHintIdx && hint.Final {
			p.stack.Cur().final = true
		}
	}
	return true, nil
}

// Handles a 'rule' parser hint
func (p *Parser) parseRule(ruleName string, sendChannel bool) (bool, error) {
	rule := grammar.GetRule(ruleName)
	//util.DumpObject(rule, "parseRule(): rule = ")
	if rule == nil {
		// TODO: make this an error once the grammar is completed
		return false, nil
	}
	p.stack.Add(rule)
	if rule.AllowFirstWordReserved {
		p.stack.Cur().allowNextWordReserved = true
	}
	e := ast.NewNode(rule.Name)
	//fmt.Printf("%#v\n", e)
	/*
		if p.stack.Cur().rule.AstFunc != nil {
			foo := p.stack.Cur().rule.AstFunc(e)
			//fmt.Printf("%s\n", foo)
			p.stack.Cur().astNode = foo
		}
	*/
	p.stack.Cur().astNode = e
	ok, err := p.parseGroup(rule.ParserHints, true)
	if err != nil {
		return false, err
	}
	curStackEntry := p.stack.Cur()
	p.stack.Remove()
	if ok {
		if p.stack.Cur() != nil {
			p.stack.Cur().astNode.AddChild(e)
		} else if sendChannel {
			util.DumpJson(e, "sending root element:\n")
			p.commandChan <- e
		}
	} else {
		if curStackEntry.final {
			token, _ := p.nextToken()
			if token != nil {
				return ok, fmt.Errorf("found unexpected token `%s` at line %d, offset %d", token.Value, token.LineNum, token.Offset)
			}
		}
	}
	return ok, nil
}

// Handles an 'any' parser hint
// Succeeds if any parser hints match
func (p *Parser) parseAny(hints []*grammar.ParserHint) (bool, error) {
	for _, hint := range hints {
		//util.DumpObject(hint, "parseAny(): hint = ")
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
	/*
		if idx >= 0 {
			fmt.Printf("setTokenIdx(%d): curToken = %#v\n", idx, p.curToken())
		}
	*/
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
