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
	p := NewParser()
	p.Start(parser_input.NewStringParserInput(input))
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
		for cmd_idx := 0; ; cmd_idx++ {
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
				t.Fatalf("more AST nodes found than expected. next node: %#v", nodes[len(outputs)])
			}
		}
	}
}

func TestParserBasic(t *testing.T) {
	/*
		test_inputs := []string{
			"foo $(echo bar foo bar) baz\nabc \"123 456\" 'd\nef' 789",
		}
		for _, input := range test_inputs {
			parser.Parse(strings.NewReader(input))
		}
	*/
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
	}
	runTests(testCases, t)
}
