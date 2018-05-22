import React from "react";
import TurtleEnableButton from "./TurtleEnableButton";
import styled from "styled-components";
import PropTypes from "prop-types";
import { buttonColor } from "../theme.js";

/**
 * Simple bar with all the turtles
 * Author: S.A. Tanja
 * Author: H.E. van der Laan
 * Author: T.T.P. Franken
 */

const TurtleEnableBar = props => {
  return (
    <div className={props.className}>
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
    </div>
  );
};

TurtleEnableBar.propTypes = {
  turtles: PropTypes.object.isRequired,
  onDisable: PropTypes.func.isRequired,
  onEnable: PropTypes.func.isRequired
};

const TurtleBar = styled(TurtleEnableBar)`
  display: flex;
  width: 100%;
  justify-content: space-between;
  overflow-x: auto;
  align-content: space-between;
  border-style: solid;
  border-width: 0rem 0rem 0.2rem 0rem;
  margin-bottom: 0.2rem;
  background: ${props => props.theme.buttonColor};
  position: fixed;
  z-index: 1;
`;

export default TurtleBar;
