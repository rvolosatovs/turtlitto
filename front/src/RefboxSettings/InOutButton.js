import React from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

/**
 * Constructs two buttons for in and out
 * Author: G.W. van der Heijden
 * Author: T.T.P. Franken
 *
 * Props:
 * - onClick: a function on what to do when the button is pressed
 */
const InOutButton = props => {
  return (
    <Container>
      <GoButton
        onClick={() => {
          props.onClick("in");
        }}
      >
        {"Go in"}
      </GoButton>
      <GoButton
        onClick={() => {
          props.onClick("out");
        }}
      >
        {"Go out"}
      </GoButton>
    </Container>
  );
};

const Container = styled.div`
  width: 50%;
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
`;

const GoButton = styled.button`
  background-color: ${props => props.theme.inOutButton};
  height: 4rem;
  width: 75%;
  margin-left: auto;
  margin-right: auto;
  border: 0.25rem solid;
  font-size: 2rem;
  display: block;
  &:active {
    background-color: ${props => props.theme.inOutButtonActive};
  }
`;

InOutButton.propType = {
  onClick: PropTypes.func.isRequired
};

export default InOutButton;
