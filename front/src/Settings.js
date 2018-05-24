import React, { Component } from "react";

import Turtle from "./Turtle";
import TurtleEnableBar from "./TurtleEnableBar";

/**
 * Show the settings of all turtles.
 * Author: H.E. van der Laan
 *
 * props:
 * - turtles: an array of Turtles
 */
export default class extends Component {
  render() {
    const { turtles } = this.props;
    return (
      <div>
        <TurtleEnableBar
          turtles={turtles}
          onDisable={position => {} /* TODO: disable turtle" */}
          onEnable={position => {} /* TODO: enable turtle */}
        />
        {turtles
          .filter(turtle => turtle.enabled)
          .map(turtle => (
            <Turtle
              key={turtle.id}
              turtle={turtle}
              editable
              onChange={(changedProp, newValue) => {} /* TODO: turtle update */}
            />
          ))}
      </div>
    );
  }
}
