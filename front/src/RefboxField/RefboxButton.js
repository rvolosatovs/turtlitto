import React from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

/**
 * A button part of the refbox. Does a callback to its parent to handle onClick events
 * Author: S.A. Tanja
 * Author: G.W. van der Heijden
 *
 * Props:
 *  - teamColor: the team color of the button (cyan or magenta).
 *  - onClick: a function on what to do when the button is pressed
 *  - children: the children of the button
 */
const RefboxButton = props => {
  return (
    <Button
      className={props.className}
      teamColor={props.teamColor}
      onClick={() => {
        props.onClick();
      }}
    >
      {props.children}
    </Button>
  );
};

const Button = styled.button`
  background-color: ${props =>
    props.teamColor === "magenta"
      ? props.theme.refboxMagenta
      : props.theme.refboxCyan};
  justify-content: center;
  align-items: center;
  font-size: 4rem;
  padding: 1.5rem 0;
  border: 0.125rem solid;
  user-select: none;
`;

RefboxButton.propType = {
  teamColor: PropTypes.oneOf(["cyan", "magenta"]).isRequired,
  tag: PropTypes.oneOf([
    "KO",
    "FK",
    "GK",
    "TI",
    "C",
    "P",
    "Soft",
    "Medium",
    "Hard"
  ]).isRequired
};

export default RefboxButton;
