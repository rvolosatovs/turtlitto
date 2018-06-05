import React from "react";
import TurtleEnableButton from "./TurtleEnableButton";
import styled from "styled-components";
import PropTypes from "prop-types";

/**
 * Simple bar with all the turtles
 * Author: S.A. Tanja
 * Author: H.E. van der Laan
 * Author: T.T.P. Franken
 * Author: B. Afonins
 *
 * Props:
 * - turtles: a list of turtles to be displayed in the bar
 *   - id: an identificator of the turtle. Should be unique.
 *   - enabled: a boolean flag specifying whether the turtle is enabled
 * - onTurtleEnableChange: function to call when the turtle enable button is pressed. The first argument is the id of the turtle that is changed.
 */

const TurtleEnableBar = props => {
  const { className, onTurtleEnableChange } = props;
  return (
    <Bar className={className}>
      {props.turtles.map(turtle => {
        return (
          <TurtleEnableButton
            key={turtle.id}
            enabled={turtle.enabled}
            id={turtle.id}
            onTurtleEnableChange={() => onTurtleEnableChange(turtle.id)}
          />
        );
      })}
    </Bar>
  );
};

TurtleEnableBar.propTypes = {
  turtles: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      enabled: PropTypes.bool.isRequired
    })
  ).isRequired,
  onTurtleEnableChange: PropTypes.func.isRequired
};

const Bar = styled.div`
  display: flex;
  height: 15%;
  width: 100%;
  justify-content: space-between;
  overflow-x: auto;
  align-content: space-between;
  border-style: solid;
  border-width: 0rem 0rem 0.2rem 0rem;
  margin-bottom: 0.2rem;
  background: ${props => props.theme.turtleEnableBar};
`;

export default TurtleEnableBar;
