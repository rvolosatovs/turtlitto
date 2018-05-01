import React from "react";
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

export default props => {
  return (
    <Input
      className={props.className}
      placeholder={props.placeholder}
      onChange={event => props.onChange(event.target.value)}
      type="text"
      value={props.value}
      disabled={props.isDisabled}
    />
  );
};
