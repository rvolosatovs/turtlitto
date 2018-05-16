import React from "react";
import TurtleEnableButton from "./TurtleEnableButton";
import styled from "styled-components";

/**
 * Simple bar with all the turtles
 * Author: S.A. Tanja
 */

const TurtleEnableBar = props => {
  return (
    <div className={props.className}>
      <TurtleEnableButton id="1" />
      <TurtleEnableButton id="2" />
      <TurtleEnableButton id="3" />
      <TurtleEnableButton id="4" />
      <TurtleEnableButton id="5" />
      <TurtleEnableButton id="6" />
    </div>
  );
};

export default TurtleEnableBar;
