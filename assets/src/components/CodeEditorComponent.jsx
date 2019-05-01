import React, { useState } from "react";
import MonacoEditor from "react-monaco-editor";
import { options, languageDef, configuration } from "./codeEditor-config.js";

export default function CodeEditorComponent() {
  const GCD = `program       GCD
  declare     x: integer, 
              y: integer
  always      decx, decy := x > y, y > x
  initially   x, y : 60, 48
  assign      x := x - y if decx
      []      y := y - x if decy
end`;
  const BubbleSort = `program       BubbleSort
  declare     N: integer, 
              a: array [0..N] of integer
  initially   N: 5 []
              <|| k: 0 <= k < N :: a[k] = rand() >
  assign      <[] i: 0 <= i < N - 1 :: a[i], a[i + 1] = a[i + 1], a[i] if a[i] > a[i + 1] >
end`;
  const [code, setCode] = useState("");

  const runCode = () => {
    const formData = new FormData();
    formData.append("code", code);
    formData.append(
      "authenticity_token",
      document.querySelector('[name="csrf-token"]').content
    );
    fetch("/backend/api/run_code/", {
      method: "POST",
      credentials: "same-origin",
      mode: "same-origin",
      headers: {
        Accept: "application/json"
      },
      body: formData
    })
      .then(response => {
        if (!response.ok) throw new Error(response.status);
        else return response.json();
      })
      .then(data => {
        data = JSON.parse(data);
        console.log(data);
        if (data.status) {
          if (window.confirm(data.result)) {
            var file_path = "/out/program.smv";
            var a = document.createElement("a");
            a.href = file_path;
            a.download = file_path.substr(file_path.lastIndexOf("/") + 1);
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
          }
        } else {
          alert(data.result);
        }
        return true;
      });
  };

  const editorDidMount = (editor, monaco) => {
    if (!monaco.languages.getLanguages().some(({ id }) => id === "unity")) {
      monaco.languages.register({ id: "unity" });
      monaco.languages.setMonarchTokensProvider("unity", languageDef);
      monaco.languages.setLanguageConfiguration("unity", configuration);
    }
    editor.focus();
  };

  return (
    <div>
      <div className="row">
        <h3>Examples: </h3>
        <span className="result_btn" href="#" onClick={() => setCode(GCD)}>
          GCD
        </span>
        <span
          className="result_btn"
          href="#"
          onClick={() => setCode(BubbleSort)}
        >
          BUBBLE SORT
        </span>
      </div>
      <br />
      <div className="row">
        <span className="result_btn" href="#" onClick={runCode}>
          RUN
        </span>
        <span className="result_btn" href="#" onClick={() => setCode("")}>
          CLEAR
        </span>
      </div>
      <MonacoEditor
        width="100%"
        height="600"
        theme="vs-dark"
        language="unity"
        value={code}
        options={options}
        onChange={(newValue, e) => setCode(newValue)}
        editorDidMount={editorDidMount}
      />
    </div>
  );
}
