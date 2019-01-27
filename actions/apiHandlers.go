package actions

import (
	"encoding/json"

	"github.com/Filip/unity_verificator/unityInterpreter"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// Result is json_result for render
type Result struct {
	Text      string                 `json:"text"`
	Variables map[string]interface{} `json:"variables"`
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
	ok := unity.Parse()
	data := Result{Text: ok, Variables: unity.Variables}
	return c.Render(jsonResponse(data))
}
