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
 *   - id: an identificator of the turtle
 *   - enabled: a boolean flag specifying whether the turtle is enabled
 * - onDisable: function to call when disabling this turtle
 * - onEnable: function to call when enabling this turtle
 */

const TurtleEnableBar = props => {
  return (
    <Bar className={props.className}>
      {props.turtles.map((turtle, position) => {
        return (
          <TurtleEnableButton
            key={turtle.id}
            enabled={turtle.enabled}
            id={position + 1}
            onEnable={() => props.onEnable(position)}
            onDisable={() => props.onDisable(position)}
          />
        );
      })}
    </Bar>
  );
};

TurtleEnableBar.propTypes = {
  turtles: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      enabled: PropTypes.bool.isRequired
    })
  ).isRequired,
  onDisable: PropTypes.func.isRequired,
  onEnable: PropTypes.func.isRequired
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
  position: sticky;
  z-index: 1;
`;

export default TurtleEnableBar;
