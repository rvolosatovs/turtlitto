import React, { Component } from "react";
import RefboxButton from "./RefboxButton";
import styled from "styled-components";

/**
 * Constructs a refbox for a team consisting of 6 buttons:
 * - Kick off (KO)
 * - Free kick (FK)
 * - Goal kick (GK)
 * - Throw in (TI)
 * - Corner (C)
 * - Penalty (P)
 * Also takes care of handeling commands that the buttons represents
 *
 * Parameter: teamColor, the color of the buttons inside of the refbox.
 *
 * Author: S.A. Tanja
 */
class RefboxField extends Component {
  constructor(props) {
    super(props);
    this.state = {
      teamColor: props.teamColor
    };
    this.HandleCommand = this.HandleCommand.bind(this);
  }

  render() {
    return (
      <div>
        <RefboxStyle>
          <RefboxButton
            color={this.state.teamColor}
            onClick={this.HandleCommand}
          >
            KO
          </RefboxButton>
          <RefboxButton
            color={this.state.teamColor}
            onClick={this.HandleCommand}
          >
            FK
          </RefboxButton>
          <RefboxButton
            color={this.state.teamColor}
            onClick={this.HandleCommand}
          >
            GK
          </RefboxButton>
          <RefboxButton
            color={this.state.teamColor}
            onClick={this.HandleCommand}
          >
            TI
          </RefboxButton>
          <RefboxButton
            color={this.state.teamColor}
            onClick={this.HandleCommand}
          >
            C
          </RefboxButton>
          <RefboxButton
            color={this.state.teamColor}
            onClick={this.HandleCommand}
          >
            P
          </RefboxButton>
        </RefboxStyle>
      </div>
    );
  }

  HandleCommand(color, command) {
    console.log("Button " + command + " was pressed from team " + color);
  }
}

/**
 * Places the buttons in 2 columns and rows of 3
 * Allows for multiple refboxes to be next to each other (float : left)
 */
const RefboxStyle = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 206px);
  grid-template-rows: repeat(3, 206px);
  margin: auto;
  width: 412px;
  float: center;
`;

export default RefboxField;
