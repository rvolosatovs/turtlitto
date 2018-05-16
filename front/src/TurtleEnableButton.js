import React from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

/**
 * A button that keeps track whether a turtle is enabled or not.
 * Is reflected to the user in terms in terms of button appearance.
 *
 * Props:
 * - enabled: whether this turtle is enabled or not
 * - onDisable: function to call when disabling this turtle
 * - onEnable: function to call when enabling this turtle
 *
 * Author: S.A. Tanja
 * Author: H.E. van der Laan
 */
const TurtleEnableButton = props => {
  if (props.enabled) {
    return (
      <Button className={props} isActive onClick={props.onDisable}>
        Turtle {props.id}
      </Button>
    );
  } else {
    return (
      <Button className={props} onClick={props.onEnable}>
        Turtle {props.id}
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
  width: 16%;
  min-width: 75px;
  height: 10vmin;
  min-height: 10%;
  font-size: 4vmin;
  min-font-size: 12px;
  flex: 1;
  margin: 1px;
`;

export default TurtleEnableButton;
