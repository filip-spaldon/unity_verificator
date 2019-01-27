// This config defines the editor's view.
export const options = {
    selectOnLineNumbers: true,
    lineNumbers: true,
    readOnly: false,
    fontSize: 12,
    tabSize: 2
};

// This config defines how the language is displayed in the editor.
export const languageDef = {
    defaultToken: '',
    number: /\d+(\.\d+)?/,
    comment: /^\s*#([ =|].*)?$/,
    out: ['program', 'end'],
    in: [
        'declare',
        'initially',
        'always',
        'assign',
        'if',
        'integer',
        'boolean',
        'min',
        'max',
        'and',
        'or',
        'true',
        'false'
    ],
    tokenizer: {
        root: [
            { include: '@whitespace' },
            { include: '@numbers' },
            { include: '@strings' },
            { include: '@tags' },
            [/^\w+/, { cases: { '@out': 'keyword' } }],
            [/\w+/, { cases: { '@in': 'keyword' } }]
        ],
        whitespace: [[/@comment/, 'comment'], [/\s+/, 'white']],
        numbers: [[/@number/, 'number']],
        strings: [
            // [/=[ @number]*$/, 'string.escape'],
            // [/:=[ @number]*$/, 'string.escape'],
            // [/[=|:][ @number|a-zA-Z]*$/, 'string.escape']
            [/[\"\'](.*?)[\"\']/, 'string.escape'],
            [/[A-Z][\w\$]*/, 'type.identifier']
        ],
        tags: [[/[+|-|=|<|>|<=|>=|==|*|/|:]/, 'tag']]
    }
};

// This config defines the editor's behavior.
export const configuration = {
    comments: {
        lineComment: '#'
    },
    brackets: [['{', '}'], ['[', ']'], ['(', ')'], ["'", "'"], ['"', '"']]
};
