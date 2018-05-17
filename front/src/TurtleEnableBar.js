import React from "react";
import TurtleEnableButton from "./TurtleEnableButton";
import styled from "styled-components";

/**
 * Simple bar with all the turtles
 * Author: S.A. Tanja
 * Author: H.E. van der Laan
 */

const TurtleEnableBar = props => {
  return (
    <div className={props.className}>
      {props.turtles.map((turtle, position) => {
        return (
          <TurtleEnableButton
            key={turtle.id}
            enabled={turtle.enabled}
            onEnable={() => props.onEnable(position)}
            onDisable={() => props.onDisable(position)}
          />
        );
      })}
    </div>
  );
};

const TurtleBar = styled(TurtleEnableBar)`
  display: flex;
  width: 100%;
  justify-content: space-between;
  overflow-x: auto;
  scrollbar: hidden;
  align-content: space-between;
  border-style: solid;
  border-width: 0px 0px 2px 0px;
  margin-bottom: 2px;
  background-color: #ededed;
  position: fixed;
  z-index: 9999;
`;

export default TurtleBar;
