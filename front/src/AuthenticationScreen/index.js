import React, { Component } from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

/**
 * An authentication screen. Upon a submit it will send the token in the
 * inputfield along with a callback function in case the token is incorrect.
 * App.js is responsible for its unrendering.
 * Author: S.A. Tanja
 *
 * Props:
 *  - onSubmit: A function that will send the token to the SRRS
 */
class AuthenticationScreen extends Component {
  constructor(props) {
    super(props);
    this.state = {
      token: "",
      showNotification: false
    };
  }

  /**
   * Callback function. In case of an incorrect token, update state to display
   * error. In case of a correct token, App.js should unrender this.
   */
  onIncorrectToken() {
    this.setState({ showNotification: true });
  }

  /**
   * Gets the user inputed token, and calls the onSubmit function given from
   * props, and gives a callback to update state.
   */
  onSubmit(event) {
    const token = this.state.token;
    this.props.onSubmit(token, isCorrectToken => this.onIncorrectToken());
  }

  render() {
    return (
      <Container>
        <Window>
          <Label>Token:</Label>
          <Input
            type="text"
            placeholder="Enter the TRC token"
            onChange={event => this.setState({ token: event.target.value })}
          />
          {this.state.showNotification && (
            <WarningLabel>Incorrect token.</WarningLabel>
          )}
          <LoginButton
            id="login-button"
            onClick={event => this.onSubmit(event)}
          >
            Log in
          </LoginButton>
        </Window>
      </Container>
    );
  }
}

//TODO: Apply theme.js for all styled components (if applicable)
const Container = styled.div`
  background: silver;
  height: 100%;
  padding-top: 10%;
`;

const Window = styled.div`
  margin: 0 auto;
  border: 0.1rem solid;
  background: #ededed;
  width: 85%;
  padding: 1rem;
  text-align: center;
`;

const Label = styled.label`
  font-size: 1.5rem;
`;

const WarningLabel = styled(Label)`
  color: red;
`;

const Input = styled.input`
  margin-top: 1rem;
  &[type="text"] {
    border: 0.1rem solid;
    width: 100%;
    padding: 1.75rem;
    font-size: 1.5rem;
  }
`;

const LoginButton = styled.button`
  margin-top: 1rem;
  background-color: white;
  border-color: black;
  border: 0.1rem solid;
  width: 100%;
  height: 5rem;
  font-size: 1.5rem;
  &:active {
    background-color: #ededed;
  }
`;

AuthenticationScreen.propTypes = {
  onSubmit: PropTypes.func.isRequired
};

export default AuthenticationScreen;
