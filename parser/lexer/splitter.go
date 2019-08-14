package lexer

import (
	//"bytes"
	"fmt"
	"github.com/agaffney/crapsh/parser/rules"
	"io"
)

func (l *Lexer) checkForOperators(value string, startsWith bool) int {
	for _, op := range rules.OperatorRules {
		var opLen int
		if startsWith {
			// Use the smaller of the value length and operator length
			opLen = len(value)
			if len(op.Pattern) < len(value) {
				opLen = len(op.Pattern)
			}
		} else {
			// Use the operator length, since we're looking for an exact match
			opLen = len(op.Pattern)
			// Move to the next operator candidate if our input is a different length than the operator
			if len(value) != opLen {
				continue
			}
		}
		// Return operator token type if there's a match
		if value[:opLen] == op.Pattern[:opLen] {
			return op.TokenType
		}
	}
	return -1
}

func (l *Lexer) checkForDelimeterStart(value string, rule *rules.DelimeterRule, startsWith bool) bool {
	var delimLen int
	if startsWith {
		// Use the smaller of the value length and the delimeter length
		delimLen = len(value)
		if len(rule.DelimStart) < len(value) {
			delimLen = len(rule.DelimStart)
		}
	} else {
		// Use the delimeter length, since we're looking for an exact match
		delimLen = len(rule.DelimStart)
		// No match if our value is smaller than our delimeter
		if len(value) < delimLen {
			return false
		}
	}
	if value[:delimLen] == rule.DelimStart[:delimLen] {
		return true
	}

	return false
}

func (l *Lexer) checkForDelimeterEnd(value string, rule *rules.DelimeterRule, startsWith bool) bool {
	var delimLen int
	if startsWith {
		// Use the smaller of the value length and the delimeter length
		delimLen = len(value)
		if len(rule.DelimEnd) < len(value) {
			delimLen = len(rule.DelimEnd)
		}
	} else {
		// Use the delimeter length, since we're looking for an exact match
		delimLen = len(rule.DelimEnd)
		// No match if our value is smaller than our delimeter
		if len(value) < delimLen {
			return false
		}
	}
	if value[:delimLen] == rule.DelimEnd[:delimLen] {
		return true
	}

	return false
}

// Tokenize the input per the POSIX spec
// https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_03
func (l *Lexer) NextToken() (*Token, error) {
	token := &Token{}
	delimRuleStack := []*rules.DelimeterRule{rules.GetDelimeterRule(`Word`)}
	processingOperator := false
	//processingDelimeter := false
	//escapeFound := false
	for {
		// Reset token line/offset if there's no value yet
		if len(token.Value) == 0 {
			token.LineNum = l.lineNum
			token.Offset = l.lineOffset
		}
		curDelimRule := delimRuleStack[len(delimRuleStack)-1]
		c, err := l.nextRune()
		if err != nil {
			if err == io.EOF {
				// Return the current token, if any
				if len(token.Value) > 0 {
					return token, err
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		if curDelimRule.AllowOperators {
			if !processingOperator {
				// Check if current rune starts an operator
				if tokenType := l.checkForOperators(string(c), true); tokenType != -1 {
					if len(token.Value) > 0 {
						// Return rune to the buffer
						l.unreadRune()
						return token, nil
					}
					processingOperator = true
				}
			}
		}
		token.Value += string(c)
		if processingOperator {
			// Check if current token value still matches an operator
			tokenType := l.checkForOperators(token.Value, false)
			if tokenType == -1 {
				// Return last rune to the buffer and return operator token
				l.unreadRune()
				token.Value = token.Value[:len(token.Value)-1]
				return token, nil
			} else {
				// Update token type with current best guess for operator
				token.Type = tokenType
			}
		}
		/*
			if !curDelimRule.IgnoreEscapes && c == '\\' {
				escapeFound = !escapeFound
				continue
			}
		*/
		ruleMatched := false
		// TODO: check for opening delimeter by start character like with operators in order to handle escapes
		for _, ruleName := range curDelimRule.AllowedRules {
			rule := rules.GetDelimeterRule(ruleName)
			if len(token.Value) >= len(rule.DelimStart) && token.Value[len(token.Value)-len(rule.DelimStart):] == rule.DelimStart {
				// Add delimeter rule to stack
				delimRuleStack = append(delimRuleStack, rule)
				ruleMatched = true
				fmt.Printf("delimeter rule start: token.Value = '%s', rule = %#v\n", token.Value, rule)
				break
			}
		}
		if ruleMatched {
			//escapeFound = false
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
				//escapeFound = true
				continue
			}
			// Remove current delimeter rule from the stack
			delimRuleStack = delimRuleStack[:len(delimRuleStack)-1]
			fmt.Printf("delimeter rule end: token.Value = '%s', rule = %#v\n", token.Value, curDelimRule)
		}
	}
	return token, nil
}
