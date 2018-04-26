import React, { Component } from "react";
import styled from "styled-components";

const Input = styled.input`
  padding: 1rem 0 1rem 1rem;
  font-size: 1.5rem;
  border-radius: 1.5rem;
  border: 0.1rem solid ${props => props.theme.border};
  transition: 0.2s ease box-shadow;

  &:hover {
    box-shadow: 0 0.3rem 0.6rem ${props => props.theme.shadow};
  }
`;

class InputField extends Component {
  state = {
    value: this.props.value,
    isDisabled: false
  };

  render() {
    return (
      <Input
        className={this.props.className}
        placeholder={this.props.placeholder}
        onChange={event => this.props.onChange(event.target.value)}
        type="text"
        value={this.state.value}
        disabled={this.state.isDisabled}
      />
    );
  }

  componentWillReceiveProps(props) {
    this.setState({
      value: props.value,
      isDisabled: props.isDisabled
    });
  }
}

export default InputField;
