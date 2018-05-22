import React from "react";
import styled from "styled-components";
import PropTypes from "prop-types";
import { buttonColor } from "../theme.js";

/**
 *
 * Author: S.A. Tanja
 * Author: H.E. van der Laan
 * Author: T.T.P. Franken
 *
 *
 * A button that keeps track whether a turtle is enabled or not.
 * Is reflected to the user in terms in terms of button appearance.
 *
 * Props:
 * - enabled: whether this turtle is enabled or not
 * - onDisable: function to call when disabling this turtle
 * - onEnable: function to call when enabling this turtle
 *
 */
const TurtleEnableButton = props => {
  if (props.enabled) {
    return (
      <Button className={props.className} isActive onClick={props.onDisable}>
        {props.id}
      </Button>
    );
  } else {
    return (
      <Button className={props} onClick={props.onEnable}>
        {props.id}
      </Button>
    );
  }
};

TurtleEnableButton.propTypes = {
  enabled: PropTypes.bool.isRequired,
  onDisable: PropTypes.func.isRequired,
  onEnable: PropTypes.func.isRequired
};

const Button = styled.button`
  border-style: ${props => (props.isActive ? "inset" : "solid")};
  border-color: ${props => (props.isActive ? "none" : buttonColor)};
  background-color: ${props => (props.isActive ? "silver" : buttonColor)};
  width: 16%;
  min-width: 7.5rem;
  height: 10rem;
  font-size: 4rem;
  flex: 1;
  margin: 0.1rem;
  &:active {
    background: "silver";
  }
`;

export default TurtleEnableButton;
