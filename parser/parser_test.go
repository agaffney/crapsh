package parser

import (
	//"fmt"
	"github.com/agaffney/crapsh/parser/ast"
	parser_input "github.com/agaffney/crapsh/parser/input"
	"github.com/agaffney/crapsh/parser/tokens"
	//"github.com/agaffney/crapsh/util"
	"io"
	"testing"
)

type parserTestCaseOutput struct {
	nodeName   string
	children   []parserTestCaseOutput
	tokenType  int
	tokenValue string
}

type parserTestCase struct {
	input   string
	outputs []parserTestCaseOutput
}

func setupParser(input string) *Parser {
	p := NewParser(parser_input.NewStringParserInput(input))
	return p
}

func getNodeChildren(node ast.Node) []ast.Node {
	var nodes []ast.Node
	nodes = append(nodes, node)
	for _, tmpNode := range node.GetChildren() {
		children := getNodeChildren(tmpNode)
		nodes = append(nodes, children...)
	}
	return nodes
}

func getTestCaseOutputs(output parserTestCaseOutput) []parserTestCaseOutput {
	var outputs []parserTestCaseOutput
	outputs = append(outputs, output)
	for _, tmpOutput := range output.children {
		children := getTestCaseOutputs(tmpOutput)
		outputs = append(outputs, children...)
	}
	return outputs
}

func runTests(testCases []parserTestCase, t *testing.T) {
	for _, testCase := range testCases {
		p := setupParser(testCase.input)
		cmd_idx := -1
		for {
			command, err := p.GetCommand()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					t.Fatalf("encountered unexpected error: %s", err.Error())
				}
			}
			if command == nil {
				break
			}
			cmd_idx++
			//util.DumpJson(command, "command = ")
			nodes := getNodeChildren(command)
			if cmd_idx >= len(testCase.outputs) {
				t.Fatalf("unexpected command")
			}
			outputs := getTestCaseOutputs(testCase.outputs[cmd_idx])
			//fmt.Printf("len(outputs) = %d, len(nodes) = %d\n", len(outputs), len(nodes))
			for idx, expected := range outputs {
				//fmt.Printf("expected[%d] = %#v\n", idx, expected)
				if idx >= len(nodes) {
					t.Fatalf("less AST nodes than expected. expecting next node '%s'", expected.nodeName)
				}
				node := nodes[idx]
				//fmt.Printf("node[%d] = %#v\n", idx, node)
				if expected.nodeName != node.GetName() {
					t.Fatalf("expected AST node type %s, found %s\n", expected.nodeName, nodes[idx].GetName())
				}
				if expected.nodeName == `Word` {
					token := node.GetToken()
					if expected.tokenType != token.Type || expected.tokenValue != token.Value {
						t.Fatalf("expected token type %d with value '%s', found token type %d with value '%s'", expected.tokenType, expected.tokenValue, token.Type, token.Value)
					}
				}
			}
			if len(outputs) < len(nodes) {
				t.Fatalf("more AST nodes found than expected. next node: %#v, next token: %#v", nodes[len(outputs)], nodes[len(outputs)].GetToken())
			}
		}
		if cmd_idx < len(testCase.outputs)-1 {
			t.Fatalf("less commands than expected")
		}
	}
}

func TestParserBasic(t *testing.T) {
	testCases := []parserTestCase{
		{
			input: `echo foo bar`,
			outputs: []parserTestCaseOutput{
				{
					nodeName: `CompleteCommand`,
					children: []parserTestCaseOutput{
						{
							nodeName: `Pipeline`,
							children: []parserTestCaseOutput{
								{
									nodeName: `SimpleCommand`,
									children: []parserTestCaseOutput{
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `echo`,
										},
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `foo`,
										},
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `bar`,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "foo $(echo bar foo bar) baz\nabc \"123 456\" 'd\nef' 789",
			outputs: []parserTestCaseOutput{
				{
					nodeName: `CompleteCommand`,
					children: []parserTestCaseOutput{
						{
							nodeName: `Pipeline`,
							children: []parserTestCaseOutput{
								{
									nodeName: `SimpleCommand`,
									children: []parserTestCaseOutput{
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `foo`,
										},
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `$(echo bar foo bar)`,
										},
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `baz`,
										},
									},
								},
							},
						},
					},
				},
				{
					nodeName: `CompleteCommand`,
					children: []parserTestCaseOutput{
						{
							nodeName: `Pipeline`,
							children: []parserTestCaseOutput{
								{
									nodeName: `SimpleCommand`,
									children: []parserTestCaseOutput{
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `abc`,
										},
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `"123 456"`,
										},
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: "'d\nef'",
										},
										{
											nodeName:   `Word`,
											tokenType:  tokens.TOKEN_WORD,
											tokenValue: `789`,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	runTests(testCases, t)
}
