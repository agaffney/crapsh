package parser

import (
	"fmt"
	"github.com/agaffney/crapsh/parser/ast"
	parser_input "github.com/agaffney/crapsh/parser/input"
	"github.com/agaffney/crapsh/parser/lexer"
	"github.com/agaffney/crapsh/parser/rules/grammar"
	//"github.com/agaffney/crapsh/util"
	"io"
)

type Parser struct {
	input    parser_input.Input
	stack    *Stack
	tokenBuf []*lexer.Token
	tokenIdx int
	lexer    *lexer.Lexer
}

func NewParser(input parser_input.Input) *Parser {
	parser := &Parser{input: input}
	parser.lexer = lexer.New(input)
	parser.Reset()
	parser.lexer.ReadLine()
	return parser
}

func (p *Parser) Reset() {
	p.stack = &Stack{parser: p}
	p.lexer.Reset()
	p.tokenBuf = make([]*lexer.Token, 0)
	p.tokenIdx = -1
}

func (p *Parser) GetCommand() (ast.Node, error) {
	// Reset the hint stack
	p.stack.Reset()
	// Start parsing with "root" element
	ok, err, commandNode := p.parseRule(`complete_command`)
	if err != nil {
		return nil, err
	}
	if !ok {
		if p.input.IsAvailable() {
			p.Reset()
			err := p.lexer.ReadLine()
			if err != nil {
				return nil, err
			}
			// Try again
			return p.GetCommand()
		} else {
			return nil, nil
		}
	}
	return commandNode, nil
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
			ok, err, _ = p.parseRule(hint.RuleName)
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
			// TODO: return "unexpected EOF" error?
			return false, nil
		}
		return false, err
	}
	//util.DumpObject(p.curToken(), "parseToken(): curToken = ")
	tokenType := p.classifyToken(token, hint)
	//fmt.Printf("classifyToken() returned type %d\n", tokenType)
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
		// TODO: provide method to roll back to TOKEN_NULL if the current rule doesn't
		// match, perhaps by passing a copy of the token instead of the original
		token.Type = tokenType
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
func (p *Parser) parseRule(ruleName string) (bool, error, ast.Node) {
	rule := grammar.GetRule(ruleName)
	//util.DumpObject(rule, "parseRule(): rule = ")
	if rule == nil {
		// TODO: make this an error once the grammar is completed
		return false, nil, nil
	}
	p.stack.Add(rule)
	if rule.AllowFirstWordReserved {
		p.stack.Cur().allowNextWordReserved = true
	}
	//fmt.Printf("%#v\n", e)
	var astNode ast.Node
	if p.stack.Cur().rule.AstFunc != nil {
		astNode = p.stack.Cur().rule.AstFunc()
	} else {
		astNode = ast.NewNode(rule.Name)
	}
	p.stack.Cur().astNode = astNode
	ok, err := p.parseGroup(rule.ParserHints, true)
	if err != nil {
		return false, err, nil
	}
	curStackEntry := p.stack.Cur()
	p.stack.Remove()
	if ok {
		// Return the node if we've reached the bottom of the parser stack
		if p.stack.Cur() == nil {
			return true, nil, astNode
		}
		p.stack.Cur().astNode.AddChild(astNode)
	} else {
		if curStackEntry.final {
			token, _ := p.nextToken()
			if token != nil {
				err := fmt.Errorf("found unexpected token `%s` at line %d, offset %d", token.Value, token.LineNum, token.Offset)
				return ok, err, nil
			}
		}
	}
	return ok, nil, nil
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
