import React from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

/**
 * A button that keeps track whether the turtle is enabled.
 * It is reflected to the user in terms in terms of button appearance.
 * Author: S.A. Tanja
 * Author: H.E. van der Laan
 * Author: T.T.P. Franken
 *
 * Props:
 *  - enabled: a boolean flag specifying whether is enabled or not
 *  - onTurtleEnableChange: function to call when the turtle enable button is pressed
 *  - id: identifier of a turtle
 *  - className: gives the classname for css
 */
const TurtleEnableButton = props => {
  const { enabled, id, className, onTurtleEnableChange } = props;

  return (
    <Button
      className={className}
      isActive={enabled}
      onClick={onTurtleEnableChange}
    >
      {id}
    </Button>
  );
};

TurtleEnableButton.propTypes = {
  enabled: PropTypes.bool.isRequired,
  onTurtleEnableChange: PropTypes.func.isRequired
};

const Button = styled.button`
  border-style: ${props =>
    props.isActive
      ? props.theme.buttonBorderStyleActive
      : props.theme.buttonBorderStyle};
  border-color: ${props =>
    props.isActive ? "none" : props.theme.buttonBorder};
  background-color: ${props =>
    props.isActive ? props.theme.buttonActive : props.theme.button};
  font-size: 4rem;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem;
  min-width: 7rem;
  flex: 1;
`;
Button.displayName = "TurtleEnableButton";

export default TurtleEnableButton;
