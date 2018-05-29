import React from "react";
import styled, { css } from "styled-components";
import PropTypes from "prop-types";

/**
 * Constructs a button for the Drop Ball command
 * Author: G.W. van der Heijden
 * Author: T.T.P. Franken
 *
 * Props:
 * - onClick: a function on what to do when the button is pressed
 */
const DropBall = props => {
  return (
    <DropBallWrapper>
      <DropBallButton
        onClick={() => {
          props.onClick();
        }}
      >
        DB
      </DropBallButton>
    </DropBallWrapper>
  );
};

const DropBallWrapper = styled.div`
  height: 25%;
  display: flex;
  flex-direction: row;
  justify-content: center;
`;

const DropBallButton = styled.button`
  background-color: ${props => props.theme.dropBallButton};
  height: 75%;
  width: 25%;
  border: 0.25rem solid;
  font-size: 2rem;
  margin: auto;
  &:active {
    background-color: ${props => props.theme.dropBallButtonActive};
  }
`;

DropBall.propType = {
  onClick: PropTypes.func.isRequired
};

export default DropBall;
