package actions

import (
	"encoding/json"

	"github.com/filip/unity_verificator/unityInterpreter"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// Tree is MarshalSerializer for struct Node
type Tree struct {
	Root      bool                     `json:"root"`
	Name      string                   `json:"name"`
	Statement string                   `json:"statement"`
	Section   string                   `json:"section"`
	Value     interface{}              `json:"value"`
	Ref       *unityInterpreter.Node   `json:"ref"`
	Nodes     []*unityInterpreter.Node `json:"nodes"`
}

// Result is json_result for render
type Result struct {
	Text        string                            `json:"result"`
	Variables   map[string]interface{}            `json:"variables"`
	Body        map[string]map[string]interface{} `json:"body"`
	ProgramName string                            `json:"Program name"`
	Tree        Tree                              `json:"tree"`
}

func jsonResponse(obj Result) (int, render.Renderer) {
	json, err := json.Marshal(obj)
	if err != nil {
		return 400, r.JSON(false)
	}
	return 200, r.JSON(string(json))
}

// runCodeAPIHandler is api for runned code
func runCodeAPIHandler(c buffalo.Context) error {
	code := c.Request().Form.Get("code")
	unity := unityInterpreter.Unity{Input: code, Index: 0, Kind: 5}
	unity.Next()
	unity.Scan()
	text, ok := unity.Parse()
	if !ok {
		data := Result{Text: text}
		return c.Render(jsonResponse(data))
	}
	data := Result{
		Text:        text,
		Variables:   unity.Variables,
		Body:        unity.Body,
		ProgramName: unity.Program,
		Tree: Tree{
			Root:      unity.Tree.Root,
			Ref:       unity.Tree.Ref,
			Nodes:     unity.Tree.Nodes,
			Name:      unity.Tree.Name,
			Statement: unity.Tree.Statement,
			Section:   unity.Tree.Section,
			Value:     unity.Tree.Value,
		},
	}
	return c.Render(jsonResponse(data))
}
