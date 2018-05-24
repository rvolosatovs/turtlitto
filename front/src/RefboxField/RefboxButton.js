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
 *  - tag: the tag of the button according to the TRC.
 */
const RefboxButton = props => {
  return (
    <Button
      teamColor={props.teamColor}
      onClick={() => {
        props.onClick();
      }}
    >
      {props.tag}
    </Button>
  );
};

/**
 * A styled component defining the style of the refbox buttons
 */
const Button = styled.button`
  background-color: ${props => props.teamColor};
  width: 20.6rem;
  height: 20.6rem;
  border: 0.25rem solid;
  font-size: 10rem;
`;

RefboxButton.propType = {
  teamColor: PropTypes.oneOf(["cyan", "magenta"]).isRequired,
  tag: PropTypes.oneOf(["KO", "FK", "GK", "TI", "C", "P"]).isRequired
};

export default RefboxButton;
