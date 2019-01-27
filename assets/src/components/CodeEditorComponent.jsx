import React from 'react';
import MonacoEditor from 'react-monaco-editor';
import { options, languageDef, configuration } from './codeEditor-config.js';

export default class CodeEditorComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      code: `program         GCD
    declare     x, y
    always      decx, decy = x > y, y > x
    initially   x, y : 10, 5
    assign      x := x - y if decx
        []      y := y - x if decy
end`
    };
    this.clearCode = this.clearCode.bind(this);
    this.onChange = this.onChange.bind(this);
    this.runCode = this.runCode.bind(this);
  }

  editorDidMount(editor, monaco) {
    if (!monaco.languages.getLanguages().some(({ id }) => id === 'unity')) {
      monaco.languages.register({ id: 'unity' });
      monaco.languages.setMonarchTokensProvider('unity', languageDef);
      monaco.languages.setLanguageConfiguration('unity', configuration);
    }
    editor.focus();
  }

  onChange(newValue, e) {
    this.setState({ code: newValue });
  }

  runCode() {
    const formData = new FormData();
    formData.append('code', this.state.code);
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
  }

  clearCode() {
    this.setState({ code: '' });
  }

  render() {
    const code = this.state.code;
    return (
      <div>
        <span className="result_btn" href="#" onClick={this.runCode}>
          Run
        </span>
        <span className="result_btn" href="#" onClick={this.clearCode}>
          Clear
        </span>
        <MonacoEditor
          width="100%"
          height="600"
          theme="vs-dark"
          language="unity"
          value={code}
          options={options}
          onChange={this.onChange}
          editorDidMount={this.editorDidMount}
        />
      </div>
    );
  }
}
