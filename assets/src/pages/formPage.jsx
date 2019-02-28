import React, { useState } from 'react';
import ReactDOM from 'react-dom';

const Form = () => {
  const form = {
    firstName: useFormInput('text', '', {
      required: 'required',
      className: 'form-control',
      id: 'id_first_name',
      name: 'first_name',
      placeholder: 'First name'
    }),
    lastName: useFormInput('text', '', {
      required: 'required',
      className: 'form-control',
      id: 'id_last_name',
      name: 'last_name',
      placeholder: 'Last name'
    }),
    age: useFormInput('number', '', {
      className: 'form-control',
      placeholder: 'Age'
    }),
    email: useFormInput('email', '', {
      className: 'form-control',
      placeholder: 'Email'
    })
  };

  const handleSubmit = e => {
    e.preventDefault();
    console.log(form);
  };

  const clearForm = () => {
    for (let key in form) {
      form[key].setValue('');
    }
  };
  return (
    <div className="row">
      <form onSubmit={handleSubmit} className="col-md-4">
        {Object.keys(form).map((key, index) => (
          <input {...form[key].props} key={index} />
        ))}
        <button type="submit">Submit</button>
        <button type="button" onClick={clearForm}>
          Clear
        </button>
      </form>
    </div>
  );
};

const useFormInput = (type = 'text', initialValue = '', attr = {}) => {
  const [value, setValue] = useState(initialValue);

  const onChange = e => {
    setValue(e.target.value);
  };
  let result = {
    props: {
      value,
      type,
      onChange,
      ...attr
    },
    setValue: setValue
  };
  return result;
};

export default Form;

ReactDOM.render(<Form />, document.getElementById('content'));
