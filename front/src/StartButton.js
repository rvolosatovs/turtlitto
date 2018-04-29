import React from "react";
import styled from "styled-components";

const Button = styled.button`
  color: ${props => (props.active ? props.theme.error : props.theme.success)};
  border: 0.2rem solid
    ${props => (props.active ? props.theme.error : props.theme.success)};
  padding: 1rem 0.5rem 1rem 0.5rem;
  font-size: 2rem;
  border-radius: 1.5rem;
  background: transparent;
  transition: color 0.5s ease, border 0.5s ease, background 0.5s ease;

  &:hover {
    background: ${props =>
      props.active ? props.theme.errorLighter : props.theme.successLighter};
  }
`;

export default props => {
  return (
    <Button
      active={props.isRunning}
      onClick={() => props.onClick()}
      className={props.className}
    >
      {props.isRunning ? "stop" : "start"}
    </Button>
  );
};
