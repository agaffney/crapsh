package lexer

import (
	"bytes"
	"github.com/agaffney/crapsh/parser/rules"
	"io"
)

func (l *Lexer) NextToken() (*Token, error) {
	token := &Token{}
	delimRuleStack := []*rules.DelimeterRule{rules.GetDelimeterRule(`Word`)}
	for {
		// Reset token line/offset if there's no value yet
		if len(token.Value) == 0 {
			token.LineNum = l.lineNum
			token.Offset = l.lineOffset
		}
		curDelimRule := delimRuleStack[len(delimRuleStack)-1]
		buf_string := l.buf.String()
		if curDelimRule.AllowOperators {
			for _, op := range rules.OperatorRules {
				if len(buf_string) >= len(op.Pattern) && buf_string[:len(op.Pattern)] == op.Pattern {
					if len(token.Value) > 0 {
						return token, nil
					}
					token.Value = op.Pattern
					token.Type = op.TokenType
					// Remove operator from buffer
					l.buf = bytes.NewBufferString(buf_string[len(op.Pattern):])
					l.lineOffset += len(op.Pattern)
					return token, nil
				}
			}
		}
		c, err := l.nextRune()
		if err != nil {
			if err == io.EOF {
				if len(token.Value) > 0 {
					return token, err
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		token.Value += string(c)
		ruleMatched := false
		for _, ruleName := range curDelimRule.AllowedRules {
			rule := rules.GetDelimeterRule(ruleName)
			if len(token.Value) >= len(rule.DelimStart) && token.Value[len(token.Value)-len(rule.DelimStart):] == rule.DelimStart {
				// Add delimeter rule to stack
				delimRuleStack = append(delimRuleStack, rule)
				ruleMatched = true
				break
			}
		}
		if ruleMatched {
			continue
		}
		if len(token.Value) >= len(curDelimRule.DelimEnd) && string(token.Value[len(token.Value)-len(curDelimRule.DelimEnd)]) == curDelimRule.DelimEnd {
			if !curDelimRule.IncludeDelim {
				if len(token.Value) > 0 {
					// Remove delimeter from token
					token.Value = token.Value[:len(token.Value)-len(curDelimRule.DelimEnd)]
				}
			}
			if curDelimRule.ReturnToken {
				if len(token.Value) > 0 {
					return token, nil
				}
				continue
			}
			// Remove current delimeter rule from the stack
			delimRuleStack = delimRuleStack[:len(delimRuleStack)-1]
		}
	}
	return token, nil
}
