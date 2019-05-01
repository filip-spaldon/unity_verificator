package unityInterpreter

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"gopkg.in/Knetic/govaluate.v3"

	"github.com/asaskevich/govalidator"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringInSliceContains(a string, list []string) bool {
	result := false
	for _, b := range list {
		if strings.Contains(a, b) {
			result = true
		}
	}
	return result
}

func forParserLeft(input string, variables map[string]interface{}) (int, string, int, string, string) {
	var i, k, N, oper1, oper2 interface{}
	Narray := []govaluate.ExpressionToken{}
	left := strings.Split(input, "::")[0]
	if expression, err := govaluate.NewEvaluableExpression(strings.Split(left, ":")[1]); err == nil {
		if expression.Tokens()[0].Kind.String() == "NUMERIC" {
			i = int(expression.Tokens()[0].Value.(float64))
		}
		if expression.Tokens()[2].Kind.String() == "VARIABLE" {
			if _, isset := variables[expression.Tokens()[2].Value.(string)]; isset {
				k = variables[expression.Tokens()[2].Value.(string)]
			} else {
				k = expression.Tokens()[2].Value
			}
		}
		if expression.Tokens()[1].Kind.String() == "COMPARATOR" {
			oper1 = expression.Tokens()[1].Value
		}
		if expression.Tokens()[3].Kind.String() == "COMPARATOR" {
			oper2 = expression.Tokens()[3].Value
		}
		Narray = expression.Tokens()[4:]
		if len(Narray) == 1 {
			if Narray[0].Kind.String() == "VARIABLE" {
				if _, isset := variables[Narray[0].Value.(string)]; isset {
					N = variables[Narray[0].Value.(string)]
				}
			} else if Narray[0].Kind.String() == "NUMERIC" {
				N = int(Narray[0].Value.(float64))
			}
		} else {
			pom := ""
			for _, s := range Narray {
				pom += fmt.Sprintf("%v", s.Value)
			}
			if expression, err := govaluate.NewEvaluableExpression(pom); err == nil {
				parameters := make(map[string]interface{}, 8)
				if len(expression.Vars()) != 0 {
					for _, val := range expression.Vars() {
						if _, isset := variables[val]; isset {
							parameters[val] = variables[val]
						}
					}
				}
				if result, err := expression.Evaluate(parameters); err == nil {
					N = int(result.(float64))
				}
			}
		}
	}
	return i.(int), k.(string), N.(int), oper1.(string), oper2.(string)
}

func forParserRightInitially(input string, variables map[string]interface{}, i int, k string, N int, oper1 string, oper2 string) ([][]int, []string) {
	functions := map[string]govaluate.ExpressionFunction{
		"rand": func(args ...interface{}) (interface{}, error) {
			return float64(rand.Intn(100)), nil
		},
	}
	var array []int
	arrays := make([][]int, 0)
	vals := make([]string, 0)
	right := strings.Split(input, "::")[1]
	_left := strings.Split(right, "=")[0]
	_right := strings.Split(right, "=")[1]
	inBrackets := false
	pomBrackets := ""
	index := 0
	for i, val1 := range strings.Split(_left, ",") {
		for _, val2 := range strings.Split(val1, " ") {
			if _, isset := variables[val2]; isset && val2 != " " {
				if strings.Contains(variables[val2].(string), "array") && strings.Contains(variables[val2].(string), "integer") {
					vals = append(vals, val2)
					array = make([]int, N)
				}
			} else if val2 == "[" {
				inBrackets = true
			} else if val2 == "]" {
				inBrackets = false
			} else if inBrackets {
				pomBrackets += val2
			}
		}
		if len(pomBrackets) > 1 {
			if expression, err := govaluate.NewEvaluableExpression(pomBrackets); err == nil && strings.Contains(pomBrackets, k) {
				parameters := make(map[string]interface{}, 8)
				parameters[k] = i
				if result, err := expression.Evaluate(parameters); err == nil {
					index = int(result.(float64))
				}
			}
		} else if pomBrackets == k {
			index = i
		}
		for {
			exp1 := strconv.Itoa(i) + oper1 + strconv.Itoa(index)
			exp2 := strconv.Itoa(index) + oper2 + strconv.Itoa(N)
			expression1, err1 := govaluate.NewEvaluableExpression(exp1)
			expression2, err2 := govaluate.NewEvaluableExpression(exp2)
			if err1 == nil && err2 == nil {
				result1, err1 := expression1.Evaluate(nil)
				result2, err2 := expression2.Evaluate(nil)
				if err1 == nil && err2 == nil {
					if result1.(bool) && result2.(bool) {
						if _rightExp, err := govaluate.NewEvaluableExpressionWithFunctions(strings.TrimSpace(_right), functions); err == nil {
							parameters := make(map[string]interface{}, 8)
							for _, val := range _rightExp.Vars() {
								if _, isset := variables[val]; isset {
									parameters[val] = variables[val]
								} else if val == k {
									parameters[val] = index
								}
							}
							if result, err := _rightExp.Evaluate(parameters); err != nil {
							} else {
								array[index] = int(result.(float64))
							}
						}
						index++
					} else {
						break
					}
				}
			}
		}
		arrays = append(arrays, array)
	}
	return arrays, vals
}

func makeAlwaysNode(body map[string]interface{}) *Node {
	pom := &Node{
		Nodes:   []*Node{},
		Section: "always",
	}
	index := 0
	for key, val := range body {
		pom.Nodes = append(pom.Nodes, &Node{
			Name: key,
		})
		if str, ok := val.(string); ok {
			for i, val := range strings.Split(str, " ") {
				if exp := makeExp(str, i, val, []string{">", "<", ">=", "<=", "==", "!="}); exp != nil {
					pom.Nodes[index].Nodes = append(pom.Nodes[index].Nodes, exp)
				}
			}
		}
		index++
	}
	return pom
}

func makeAssignNode(body map[string]interface{}, Tree Node) *Node {
	pom := &Node{Nodes: []*Node{}, Section: "assign"}
	index := 0
	for key, val := range body {
		if !strings.Contains(key, "for") {
			if str, ok := val.(string); ok {
				if strings.Contains(str, "if") {
					r := strings.Split(str, "if ")[1]
					str = strings.TrimSpace(strings.Split(str, "if ")[0])
					if len(strings.Split(r, " ")) == 1 {
						ref := strings.Split(r, " ")[0]
						pom.Nodes = append(pom.Nodes, &Node{Statement: "=", Section: "subAssign", Ref: Tree.Find(ref).Nodes[0]})
					} else {
						for i, val := range strings.Split(r, " ") {
							if exp := makeExp(r, i, val, []string{">", "<", ">=", "<=", "==", "!="}); exp != nil {
								pom.Nodes = append(pom.Nodes, &Node{Statement: "=", Ref: exp, Section: "subAssign"})
							}
						}
					}
				} else {
					pom.Nodes = append(pom.Nodes, &Node{Statement: "=", Section: "subAssign"})
				}
				pom.Nodes[index].Nodes = append(pom.Nodes[index].Nodes, &Node{Name: key})
				pom.Nodes[index].Nodes = append(pom.Nodes[index].Nodes, &Node{Value: str})
			}
			index++
		} else {
			forNode := &Node{Nodes: []*Node{}, Statement: "for", Section: "subAssign", Value: key[strings.Index(key, "for")+5:]}
			if str, ok := val.(string); ok {
				if strings.Contains(str, "if") {
					r := strings.Split(str, "if ")[1]
					str = strings.Split(str, "if ")[0]
					if len(strings.Split(r, " ")) == 1 {
						ref := strings.Split(r, " ")[0]
						forNode.Nodes = append(forNode.Nodes, &Node{Statement: "=", Ref: Tree.Find(ref).Nodes[0]})
					} else {
						_l := ""
						_r := ""
						_oper := ""
						afterOper := false
						for _, val := range strings.Split(r, " ") {
							if !afterOper && !stringInSlice(val, []string{">", "<", ">=", "<=", "==", "!="}) {
								_l += val
							} else if stringInSlice(val, []string{">", "<", ">=", "<=", "==", "!="}) {
								_oper += val
								afterOper = true
							} else if afterOper && !stringInSlice(val, []string{">", "<", ">=", "<=", "==", "!="}) {
								_r += val
							}
						}
						exp := &Node{Statement: _oper, Nodes: []*Node{&Node{Value: _l}, &Node{Value: _r}}}
						forNode.Nodes = append(forNode.Nodes, &Node{Statement: "=", Ref: exp})
					}
				} else {
					forNode.Nodes = append(forNode.Nodes, &Node{Statement: "="})
				}
				forNode.Nodes[0].Nodes = append(forNode.Nodes[0].Nodes, &Node{Statement: "arrExp", Value: strings.TrimSpace(str)})
			}
			pom.Nodes = append(pom.Nodes, forNode)
		}
	}
	return pom
}

func makeExp(str string, i int, val string, exps []string) *Node {
	if !govalidator.IsInt(val) && !govalidator.IsAlpha(val) && stringInSlice(val, exps) {
		exp := &Node{Statement: val, Nodes: []*Node{}}
		l := strings.Split(str, " ")[i-1]
		r := strings.Split(str, " ")[i+1]
		if govalidator.IsInt(l) {
			exp.Nodes = append(exp.Nodes, &Node{Value: l})
		} else if govalidator.IsAlpha(l) {
			exp.Nodes = append(exp.Nodes, &Node{Name: l})
		}
		if govalidator.IsInt(r) {
			exp.Nodes = append(exp.Nodes, &Node{Value: r})
		} else if govalidator.IsAlpha(r) {
			exp.Nodes = append(exp.Nodes, &Node{Name: r})
		}
		return exp
	}
	return nil
}

func getIndexArray(str string, index int) string {
	inBrackets := false
	pomBrackets := ""
	value := ""
	for _, val := range str {
		val := string(val)
		if val == "[" {
			inBrackets = true
		} else if val == "]" {
			inBrackets = false
		} else if inBrackets {
			pomBrackets += val
		} else if val != " " {
			value += val
		}
	}
	if expression, err := govaluate.NewEvaluableExpression(pomBrackets); err == nil {
		parameters := make(map[string]interface{}, 8)
		for _, val := range expression.Vars() {
			parameters[val] = index
		}
		if result, err := expression.Evaluate(parameters); err == nil {
			return fmt.Sprintf("%s[%d]", value, int(result.(float64)))
		}
	}
	return ""
}

func makeSwap(left string, right string) string {
	pom := ""
	pom += fmt.Sprintf("tmp = %s;\n", left)
	pom += fmt.Sprintf("%s = %s;\n", left, right)
	pom += fmt.Sprintf("%s = tmp;\n", right)
	return pom
}

// Node sflkajsdf
type Node struct {
	Root      bool
	Name      string
	Statement string
	Section   string
	Value     interface{}
	Ref       *Node
	Nodes     []*Node
}

func (n *Node) Find(name string) *Node {
	queue := make([]*Node, 0)
	queue = append(queue, n)
	for len(queue) > 0 {
		nextUp := queue[0]
		queue = queue[1:]
		if nextUp.Name == name {
			return nextUp
		}
		if len(nextUp.Nodes) > 0 {
			for _, child := range nextUp.Nodes {
				queue = append(queue, child)
			}
		}
	}
	return nil
}
