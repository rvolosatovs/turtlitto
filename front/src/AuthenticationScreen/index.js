import React, { Component } from "react";
import styled from "styled-components";
import PropTypes from "prop-types";

//Move these two to a different location?
import ConnectionBar from "../BottomBar/ConnectionBar";

/**
 * An authentication screen. Upon a submit it will send the token in the
 * inputfield along with a callback function in case the token is incorrect.
 * App.js is responsible for its unrendering.
 * Author: S.A. Tanja
 *
 * Props:
 *  - onSubmit: A function that will send the token to the SRRS
 *  - connectionStatus: The current connection status
 */

class AuthenticationScreen extends Component {
  constructor(props) {
    super(props);
    this.state = {
      token: ""
    };
  }

  render() {
    const { connectionStatus, notification } = this.props;

    return (
      <Container>
        <Window>
          <Label>Token:</Label>
          <Input
            type="text"
            placeholder="Enter the TRC token"
            onChange={event => this.setState({ token: event.target.value })}
          />
          {notification && (
            <WarningLabel id="auth-screen__warning-label">
              {notification}
            </WarningLabel>
          )}
          <LoginButton
            id="login-button"
            onClick={() => this.props.onSubmit(this.state.token)}
          >
            Log in
          </LoginButton>
        </Window>
        <ConnectionWindow>
          <ConnectionBar connectionStatus={connectionStatus} />
        </ConnectionWindow>
      </Container>
    );
  }
}

const Container = styled.div`
  background: ${props => props.theme.loginScreenBackground};
  height: 100%;
  padding-top: 10%;
`;

const Window = styled.div`
  margin: 0 auto;
  border: 0.1rem solid;
  border-bottom: none;
  background: ${props => props.theme.loginFormBackground};
  width: 85%;
  padding: 1rem;
  text-align: center;
`;

const Label = styled.label`
  font-size: 1.5rem;
`;

const WarningLabel = styled(Label)`
  color: ${props => props.theme.error};
`;
WarningLabel.displayName = "WarningLabel";

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
  background-color: ${props => props.theme.loginButton};
  border: 0.1rem solid;
  width: 100%;
  height: 5rem;
  font-size: 1.5rem;

  &:active {
    background-color: ${props => props.theme.loginButtonActive};
  }
`;

const ConnectionWindow = styled(Window)`
  padding: 0;
  border-bottom: 0.1rem solid;
`;

AuthenticationScreen.propTypes = {
  onSubmit: PropTypes.func.isRequired,
  notification: PropTypes.string.isRequired
};

export default AuthenticationScreen;
