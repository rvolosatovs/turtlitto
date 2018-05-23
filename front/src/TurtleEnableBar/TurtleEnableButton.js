import React from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

const onClick = (enabled, onEnable, onDisable) => {
  if (enabled) {
    onDisable();
  } else {
    onEnable();
  }
};
/**
 * A button that keeps track whether the turtle is enabled
 * Is reflected to the user in terms in terms of button appearance
 * Author: S.A. Tanja
 * Author: H.E. van der Laan
 * Author: T.T.P. Franken
 *
 * Props:
 * - enabled: a boolean flag specifying whether is enabled or not
 * - onDisable: function to call when disabling this turtle
 * - onEnable: function to call when enabling this turtle
 * - id: identificator of a turtle
 */
const TurtleEnableButton = props => {
  const { enabled, onDisable, onEnable, id, className } = props;

  return (
    <Button
      className={className}
      isActive={enabled}
      onClick={() => onClick(enabled, onEnable, onDisable)}
    >
      {id}
    </Button>
  );
};

TurtleEnableButton.propTypes = {
  enabled: PropTypes.bool.isRequired,
  onDisable: PropTypes.func.isRequired,
  onEnable: PropTypes.func.isRequired
};

const Button = styled.button`
  border-style: ${props => (props.isActive ? "inset" : "solid")};
  border-color: ${props =>
    props.isActive ? "none" : props.theme.turtleEnableButton};
  background-color: ${props =>
    props.isActive
      ? props.theme.turtleEnableButtonActive
      : props.theme.turtleEnableButton};
  width: 16%;
  min-width: 7.5rem;
  height: 10rem;
  font-size: 4rem;
  flex: 1;
  margin: 0.1rem;
`;

export default TurtleEnableButton;
