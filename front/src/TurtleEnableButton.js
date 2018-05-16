import React, { Component } from "react";
import styled from "styled-components";

/**
 * A button that keeps track whether a turtle is enabled or not.
 * Is reflected to the user in terms in terms of button appearance.
 * Author: S.A. Tanja
 */
class TurtleEnableButton extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isActive: false,
      id: props.id
    };
    this.HandleClick = this.HandleClick.bind(this);
  }

  render() {
    return (
      <Button isActive={this.state.isActive} onClick={this.HandleClick}>
        Turtle {this.state.id}
      </Button>
    );
  }

  HandleClick(props) {
    if (!this.state.isActive) {
      this.setState({ isActive: true });
      console.log("Enabling turtle " + this.state.id + "...");
    } else {
      this.setState({ isActive: false });
      console.log("Disabling turtle " + this.state.id + "...");
    }
  }
}

const Button = styled.button`
  border-style: ${props => (props.isActive ? "inset" : "solid")};
  width: 207px;
  height: 166px;
  font-size: 40px;
`;

export default TurtleEnableButton;
