import React from "react";
import TurtleEnableButton from "./TurtleEnableButton";

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

export default TurtleEnableBar;
