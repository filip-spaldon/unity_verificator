import React, { useState, useEffect } from 'react';
import MonacoEditor from 'react-monaco-editor';
import { options, languageDef, configuration } from './codeEditor-config.js';

export default function CodeEditorComponent() {
  const [code, setCode] = useState(
    `program         GCD
    declare     x: integer, 
                y: integer
    always      decx, decy = x > y, y > x
    initially   x, y : 10, 5
    assign      x := x - y if decx
        []      y := y - x if decy
end`
  );

  const runCode = () => {
    const formData = new FormData();
    formData.append('code', code);
    formData.append(
      'authenticity_token',
      document.querySelector('[name="csrf-token"]').content
    );
    fetch('/backend/api/run_code/', {
      method: 'POST',
      credentials: 'same-origin',
      mode: 'same-origin',
      headers: {
        Accept: 'application/json'
      },
      body: formData
    })
      .then(response => {
        if (!response.ok) throw new Error(response.status);
        else return response.json();
      })
      .then(data => {
        console.log(JSON.parse(data));
        return true;
      });
  };

  const editorDidMount = (editor, monaco) => {
    if (!monaco.languages.getLanguages().some(({ id }) => id === 'unity')) {
      monaco.languages.register({ id: 'unity' });
      monaco.languages.setMonarchTokensProvider('unity', languageDef);
      monaco.languages.setLanguageConfiguration('unity', configuration);
    }
    editor.focus();
  };

  return (
    <div>
      <span className="result_btn" href="#" onClick={runCode}>
        Run
      </span>
      <span className="result_btn" href="#" onClick={() => setCode('')}>
        Clear
      </span>
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
