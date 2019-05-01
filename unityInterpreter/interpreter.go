package unityInterpreter

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"gopkg.in/Knetic/govaluate.v3"
)

// inicialization vars
var (
	NOTHING     = 0
	NUMBER      = 1
	WORD        = 2
	SYMBOL      = 3
	BULDINWORDS = []string{"if", "program", "end", "declare", "always", "initially", "assign", "boolean", "integer", "true", "false", "min", "max", "and", "or", "array"}
)

// Unity - Verification tool
type Unity struct {
	Program   string
	Input     string
	Index     int
	Look      string
	Token     string
	Kind      int
	Position  int
	Variables map[string]interface{}
	Body      map[string]map[string]interface{}
	Tree      Node
}

func (u *Unity) Next() {
	if u.Index >= len(u.Input) {
		u.Look = string(0)
	} else {
		u.Look = string(u.Input[u.Index])
		u.Index++
	}
}

func (u *Unity) Scan() {
	for u.Look == " " || u.Look == "\n" {
		u.Next()
	}
	u.Token = ""
	u.Position = u.Index - 1
	if govalidator.IsInt(u.Look) {
		for govalidator.IsInt(u.Look) {
			u.Token += u.Look
			u.Next()
		}
		u.Kind = NUMBER
	} else if govalidator.IsAlpha(u.Look) {
		for govalidator.IsAlpha(u.Look) {
			u.Token += u.Look
			u.Next()
		}
		u.Kind = WORD
	} else if !govalidator.IsAlpha(u.Look) && !govalidator.IsInt(u.Look) && u.Look != string(0) && u.Look != " " && u.Look != "\n" {
		for !govalidator.IsAlpha(u.Look) && !govalidator.IsInt(u.Look) && u.Look != string(0) && u.Look != " " && u.Look != "\n" {
			u.Token += u.Look
			u.Next()
		}
		u.Kind = SYMBOL
	} else if u.Look == string(0) {
		u.Kind = NOTHING
	}
}

func (u *Unity) Parse() (string, bool) {
	u.Variables = make(map[string]interface{})
	u.Body = make(map[string]map[string]interface{})
	u.Tree = Node{
		Nodes: []*Node{},
		Root:  true,
	}
	for u.Kind != NOTHING {
		if u.Token == "program" {
			u.Scan()
			u.Program = u.Token
			u.Tree.Name = u.Program
		} else if u.Token == "declare" {
			u.Body["declare"] = make(map[string]interface{})
			declareNode := &Node{
				Nodes:   []*Node{},
				Section: "declare",
			}
			u.Scan()
			pom := ""
			for u.Token != "always" && u.Token != "initially" {
				pom += u.Token + " "
				u.Scan()
			}
			for _, val := range strings.Split(pom, ",") {
				left := strings.TrimSpace(strings.Split(val, ":")[0])
				right := strings.Split(val, ":")[1]
				right = right[1 : len(right)-1]
				if stringInSliceContains(right, BULDINWORDS) {
					u.Variables[left] = right
					u.Body["declare"][left] = right
					declareNode.Nodes = append(declareNode.Nodes, &Node{
						Name:      left,
						Statement: right,
					})
				} else {
					return "CHYBA - declare section", false
				}
			}
			u.Tree.Nodes = append(u.Tree.Nodes, declareNode)
		} else if u.Token == "always" {
			u.Body["always"] = make(map[string]interface{})
			u.Scan()
			pom := ""
			for u.Token != "initially" {
				pom += u.Token + " "
				u.Scan()
			}
			left := strings.Split(pom, ":=")[0]
			right := strings.Split(pom, ":=")[1]
			if len(strings.Split(left, ",")) != len(strings.Split(right, ",")) {
				return "CHYBA - always section", false
			}
			for index, val := range strings.Split(left, ",") {
				rightVal := string(strings.Split(right, ",")[index])
				val := strings.TrimSpace(val)
				u.Variables[val] = rightVal[1 : len(rightVal)-1]
				u.Body["always"][val] = rightVal[1 : len(rightVal)-1]
			}
			alwaysNode := makeAlwaysNode(u.Body["always"])
			u.Tree.Nodes = append(u.Tree.Nodes, alwaysNode)
		} else if u.Token == "initially" {
			u.Body["initially"] = make(map[string]interface{})
			initiallyNode := &Node{
				Nodes:   []*Node{},
				Section: "initially",
			}
			u.Scan()
			pom := ""
			for u.Token != "assign" {
				pom += u.Token + " "
				u.Scan()
			}
			for _, val := range strings.Split(pom, " [] ") {
				if strings.Contains(val, "<[]") || strings.Contains(val, "<||") {
					node := &Node{
						Statement: "=",
						Nodes:     []*Node{},
					}
					val = strings.TrimSpace(val)
					val = val[strings.Index(val, "<")+4 : len(val)-1]
					i, k, N, oper1, oper2 := forParserLeft(val, u.Variables)
					arrays, vals := forParserRightInitially(val, u.Variables, i, k, N, oper1, oper2)
					for index, val := range vals {
						u.Variables[val] = arrays[index]
						u.Body["initially"][val] = arrays[index]
						node.Nodes = append(node.Nodes, &Node{Name: val})
						node.Nodes = append(node.Nodes, &Node{Value: arrays[index]})
					}
					initiallyNode.Nodes = append(initiallyNode.Nodes, node)
				} else {
					left := strings.TrimSpace(strings.Split(val, ":")[0])
					right := strings.TrimSpace(strings.Split(val, ":")[1])
					if len(strings.Split(left, ",")) != len(strings.Split(right, ",")) {
						return "CHYBA - initially section", false
					}
					for index, _val := range strings.Split(left, ",") {
						node := &Node{
							Statement: "=",
							Nodes:     []*Node{},
						}
						_val = strings.TrimSpace(_val)
						_r := strings.TrimSpace(strings.Split(right, ",")[index])
						if _, isset := u.Variables[_val]; isset {
							if a, err := strconv.Atoi(_r); err == nil && u.Variables[_val] == "integer" {
								u.Variables[_val] = a
								u.Body["initially"][_val] = a
								node.Nodes = append(node.Nodes, &Node{Name: _val})
								node.Nodes = append(node.Nodes, &Node{Value: a})
							}
						} else {
							return "CHYBA - initially section", true
						}
						initiallyNode.Nodes = append(initiallyNode.Nodes, node)
					}
				}
			}
			u.Tree.Nodes = append(u.Tree.Nodes, initiallyNode)
		} else if u.Token == "assign" {
			u.Body["assign"] = make(map[string]interface{})
			u.Scan()
			pom := ""
			for u.Token != "end" {
				pom += u.Token + " "
				u.Scan()
			}
			for index, val := range strings.Split(pom, " [] ") {
				if strings.Contains(val, "<[]") || strings.Contains(val, "<||") {
					val = strings.TrimSpace(val)
					val = val[strings.Index(val, "<")+4 : len(val)-1]
					left := strings.TrimSpace(strings.Split(val, "::")[0])
					right := strings.TrimSpace(strings.Split(val, "::")[1])
					u.Body["assign"]["for_"+strconv.Itoa(index)+" "+strings.TrimSpace(left)] = strings.TrimSpace(right)
				} else {
					val = strings.TrimSpace(val)
					left := strings.TrimSpace(strings.Split(val, ":=")[0])
					right := strings.TrimSpace(strings.Split(val, ":=")[1])
					u.Body["assign"][left] = right

				}
			}
			assignNode := makeAssignNode(u.Body["assign"], u.Tree)
			u.Tree.Nodes = append(u.Tree.Nodes, assignNode)
		} else if u.Token == "end" {
			return "Program bol sparsovaný správne", true
		} else {
			u.Scan()
		}
	}
	return "Telo programu je prázdne", false
}

func MakePromela(root *Node, u *Unity) {
	var (
		_declare   = true
		_initially = false
		_always    = false
		_assign    = false
	)
	queue := make([]*Node, 0)
	queue = append(queue, root.Nodes...)
	assignIndex := 0
	var pom, declarePom, initiallyPom string
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]
		if _declare {
			if current.Section == "always" {
				_declare = false
				_always = true
			} else if current.Section == "initially" {
				_declare = false
				_initially = true
			} else if current.Section == "" {
			}
		} else if _always {
			if current.Section == "initially" {
				_always = false
				_initially = true
			} else {
			}
		} else if _initially {
			if current.Section == "assign" {
				_initially = false
				_assign = true
				pom += "int tmp;\n"
				pom += declarePom
				pom += "init {\n"
				pom += initiallyPom
				pom += "}\n"
			} else {
				if current.Statement == "=" && len(current.Nodes) == 2 {
					switch x := current.Nodes[1].Value.(type) {
					case int:
						declarePom += fmt.Sprintf("int %s;\n", current.Nodes[0].Name)
						initiallyPom += fmt.Sprintf("%s = %d;\n", current.Nodes[0].Name, x)
					case []int:
						declarePom += fmt.Sprintf("int %s[%d];\n", current.Nodes[0].Name, len(x))
						for index, val := range x {
							initiallyPom += fmt.Sprintf("%s[%d] = %d;\n", current.Nodes[0].Name, index, val)
						}
					}
				}
			}
		} else if _assign {
			if current.Statement == "for" && current.Section == "subAssign" {
				i, _, N, oper1, oper2 := forParserLeft(current.Value.(string), u.Variables)
				index := i
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
								pom += fmt.Sprint("active proctype process_", assignIndex, "() {\n")
								pom += "do\n"
								if current.Nodes[0].Statement == "=" && current.Nodes[0].Ref != nil {
									l := getIndexArray(current.Nodes[0].Ref.Nodes[0].Value.(string), index)
									r := getIndexArray(current.Nodes[0].Ref.Nodes[1].Value.(string), index)
									pom += fmt.Sprintf(":: %s %s %s ->\n", l, current.Nodes[0].Ref.Statement, r)
								} else {
									pom += ":: "
								}
								pom += "atomic {\n"
								if current.Nodes[0].Nodes[0].Statement == "arrExp" {
									value := current.Nodes[0].Nodes[0].Value.(string)
									l := strings.Split(value, "=")[0]
									r := strings.Split(value, "=")[1]
									lArray := make([]string, 0)
									rArray := make([]string, 0)
									for i, val := range strings.Split(l, ",") {
										lArray = append(lArray, strings.TrimSpace(val))
										rArray = append(rArray, strings.TrimSpace(strings.Split(r, ",")[i]))
									}
									if lArray[0] == rArray[1] {
										pom += makeSwap(getIndexArray(lArray[0], index), getIndexArray(lArray[1], index))
									}
								}
								pom += "}\n"
								pom += ":: else -> skip\n"
								pom += "od\n"
								pom += "}\n"
								assignIndex++
								index++
							} else {
								break
							}
						}
					}
				}
			} else if current.Statement == "=" && current.Section == "subAssign" {
				pom += fmt.Sprint("active proctype process_", assignIndex, "() {\n")
				pom += "do\n"
				if current.Ref != nil {
					pom += ":: "
					if current.Ref.Nodes[0].Name != "" {
						pom += fmt.Sprint(current.Ref.Nodes[0].Name, " ")
					} else if current.Ref.Nodes[0].Value != nil {
						pom += fmt.Sprint(current.Ref.Nodes[0].Value, " ")
					}
					pom += fmt.Sprint(current.Ref.Statement, " ")
					if current.Ref.Nodes[1].Name != "" {
						pom += fmt.Sprint(current.Ref.Nodes[1].Name, " ->\n")
					} else if current.Ref.Nodes[1].Value != nil {
						pom += fmt.Sprint(current.Ref.Nodes[1].Value, " ->\n")
					}
				}
				pom += "atomic {\n"
				pom += fmt.Sprintln(current.Nodes[0].Name, " = ", current.Nodes[1].Value)
				pom += "}\n"
				pom += ":: else -> skip\n"
				pom += "od\n"
				pom += "}\n"
				assignIndex++
			}
		}
		for i := range current.Nodes {
			queue = append([]*Node{current.Nodes[len(current.Nodes)-1-i]}, queue...)
		}
	}
	os.Mkdir("public/out", 0777)
	if f, err := os.Create("public/out/program.pml"); err == nil {
		if _, err := f.WriteString(pom); err == nil {
			err = f.Close()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
