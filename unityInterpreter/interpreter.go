package unityInterpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
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

func (u *Unity) Check(kind int) string {
	switch kind {
	case 0:
		return "NOTHING"
	case 1:
		return "NUMBER"
	case 2:
		return "WORD"
	case 3:
		return "SYMBOL"
	default:
		return ""
	}
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
			if len(strings.Split(pom, "[]")) <= 1 {
				left := strings.TrimSpace(strings.Split(pom, ":")[0])
				right := strings.Split(pom, ":")[1]
				right = right[1 : len(right)-1]
				if len(strings.Split(left, ",")) != len(strings.Split(right, ",")) {
					return "CHYBA - initially section", false
				}
				for index, val := range strings.Split(left, ",") {
					val = strings.TrimSpace(val)
					node := &Node{
						Statement: "=",
						Nodes:     []*Node{},
					}
					if _, isset := u.Variables[val]; isset {
						value := strings.TrimSpace(string(strings.Split(right, ",")[index]))
						if a, err := strconv.Atoi(value); err == nil && u.Variables[val] == "integer" {
							u.Variables[val] = a
							u.Body["initially"][val] = a
							node.Nodes = append(node.Nodes, &Node{Name: val})
							node.Nodes = append(node.Nodes, &Node{Value: a})
						}
					} else {
						return "CHYBA - initially section", false
					}
					initiallyNode.Nodes = append(initiallyNode.Nodes, node)
				}
				u.Tree.Nodes = append(u.Tree.Nodes, initiallyNode)
			} else {
				for _, val := range strings.Split(pom, "[]") {
					node := &Node{
						Statement: "=",
						Nodes:     []*Node{},
					}
					if !strings.Contains(val, "<<") && !strings.Contains(val, ">>") {
						left := strings.TrimSpace(strings.Split(val, ":")[0])
						right := strings.TrimSpace(strings.Split(val, ":")[1])
						if _, isset := u.Variables[left]; isset {
							if a, err := strconv.Atoi(right); err == nil && u.Variables[left] == "integer" {
								u.Variables[left] = a
								u.Body["initially"][left] = a
								node.Nodes = append(node.Nodes, &Node{Name: left})
								node.Nodes = append(node.Nodes, &Node{Value: a})
							}
						} else {
							return "CHYBA - initially section", false
						}
					} else {
						i, k, N, oper1, oper2 := forParserLeft(val[1:len(val)-1], u.Variables)
						arrays, vals := forParserRight(val[1:len(val)-1], u.Variables, i, k, N, oper1, oper2)
						for index, val := range vals {
							u.Variables[val] = arrays[index]
							u.Body["initially"][val] = arrays[index]
							node.Nodes = append(node.Nodes, &Node{Name: val})
							node.Nodes = append(node.Nodes, &Node{Value: arrays[index]})
						}
					}
					initiallyNode.Nodes = append(initiallyNode.Nodes, node)
				}
				u.Tree.Nodes = append(u.Tree.Nodes, initiallyNode)
			}
		} else if u.Token == "assign" {
			u.Body["assign"] = make(map[string]interface{})
			u.Scan()
			pom := ""
			for u.Token != "end" {
				pom += u.Token + " "
				u.Scan()
			}
			if strings.Contains(pom, "<<[]") {
				pom = pom[strings.Index(pom, "<<[]")+4:strings.Index(pom, ">>")] + pom[strings.Index(pom, ">>")+2:]
				for _, val := range strings.Split(pom, "[]") {
					fmt.Println(val)
				}
			} else if strings.Contains(pom, "<<||") {
				pom = pom[strings.Index(pom, "<<||")+4:strings.Index(pom, ">>")] + pom[strings.Index(pom, ">>")+2:]
				for _, val := range strings.Split(pom, "[]") {
					fmt.Println(val)
					//TODO ked bude cas :D
				}
			} else {
				for _, val := range strings.Split(pom, "[]") {
					left := strings.Split(val, ":=")[0]
					right := strings.Split(val, ":=")[1]
					u.Body["assign"][strings.TrimSpace(left)] = right[1 : len(right)-1]
				}
				assignNode := makeAssignNode(u.Body["assign"], u.Tree)
				u.Tree.Nodes = append(u.Tree.Nodes, assignNode)
			}
		} else if u.Token == "end" {
			return "Program bol sparsovaný správne", true
		} else {
			u.Scan()
		}
	}
	return "Telo programu je prázdne", false
}
