package unityInterpreter

import (
	"strings"

	"github.com/asaskevich/govalidator"
)

// inicialization vars
var (
	NOTHING     = 0
	NUMBER      = 1
	WORD        = 2
	SYMBOL      = 3
	BULDINWORDS = []string{"if", "program", "end", "declare", "always", "initially", "assign", "boolean", "integer", "true", "false", "min", "max", "and", "or"}
)

// Logo - Verification tool
type Unity struct {
	Program   string
	Input     string
	Index     int
	Look      string
	Token     string
	Kind      int
	Position  int
	Variables map[string]interface{}
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
	} else if !govalidator.IsAlpha(u.Look) && !govalidator.IsInt(u.Look) && u.Look != string(0) {
		for !govalidator.IsAlpha(u.Look) && !govalidator.IsInt(u.Look) {
			u.Token += u.Look
			u.Next()
		}
		u.Kind = SYMBOL
	} else if u.Look == string(0) {
		u.Kind = NOTHING
	}
}

func (u *Unity) Parse() string {
	u.Variables = make(map[string]interface{})
	for u.Kind != NOTHING {
		if u.Token == "program" {
			u.Scan()
			u.Program = u.Token
		} else if u.Token == "declare" {
			u.Scan()
			for u.Token != "always" && u.Token != "initially" {
				if u.Kind == WORD {
					u.Variables[u.Token] = nil
				}
				u.Scan()
			}
		} else if u.Token == "always" {
			u.Scan()
			pom := ""
			for u.Token != "initially" {
				token := strings.TrimSpace(u.Token)
				pom += token
				u.Scan()
			}
			left := strings.Split(pom, "=")[0]
			right := strings.Split(pom, "=")[1]
			if len(strings.Split(left, ",")) != len(strings.Split(right, ",")) {
				return "CHYBA, always section"
			}
			for index, val := range strings.Split(left, ",") {
				u.Variables[val] = string(strings.Split(right, ",")[index])
			}
		} else if u.Token == "initially" {
			u.Scan()
			pom := ""
			for u.Token != "assign" {
				token := strings.TrimSpace(u.Token)
				pom += token
				u.Scan()
			}
			left := strings.Split(pom, ":")[0]
			right := strings.Split(pom, ":")[1]
			for index, val := range strings.Split(left, ",") {
				if _, isset := u.Variables[val]; isset {
					u.Variables[val] = string(strings.Split(right, ",")[index])
				} else {
					return "CHYBA, initially section"
				}
			}
		} else if u.Token == "end" {
			return "Program bol sparsovaný správne"
		} else {
			u.Scan()
		}
	}
	return "Telo programu je prázdne"
}

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	log.Printf("%s took %fs", name, elapsed.Seconds())
// }

// func main() {
// 	defer timeTrack(time.Now(), "Logo time")
// 	logo := Logo{input: "opakuj 4 [dopredu 100 vpravo 90]", Index: 0, Kind: 5}
// 	logo.Next()
// 	logo.scan()
// 	logo.interpret()
// }
