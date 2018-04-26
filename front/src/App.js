import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import "flexboxgrid2/flexboxgrid2.css";
import styled, { css } from "styled-components";

import StartButton from "./StartButton";
import NumberInput from "./NumberInput";

const Header = styled.header`
  height: 100vh;
  background: linear-gradient(270deg, #196ebd, #01b0dd);
`;

const Content = styled.div`
  margin: 0;
  display: flex;
  flex-direction: column;
  align-items: space-around;
  justify-content: center;
  height: 100%;
`;

const Title = styled.h1`
  font-size: 3.6rem;
  color: white;
`;

const MessagesList = styled.ul`
  list-style-type: none;
  height: 40rem; // dont do that
  margin: 0 auto;
  border: 0.2rem solid #d3d3d3;
  background: white;
  overflow-y: scroll;
  overflow-x: hidden;
  text-align: left;
  padding: 0.6rem;
  width: 80%;
`;

const Message = styled.li`
  font-size: 1.5rem;
  padding: 0.5rem 0;
`;

const InfoMessage = styled.li`
  color: #bdbdbd;
`;

const ErrorMessage = Message.extend`
  color: red;
`;

const PingButton = styled.button`
  padding: 0.5rem 1rem;
  border-radius: 1.5rem;
  text-transform: uppercase;
  border: 0.2rem solid ${props => (props.isDisabled ? "rgba(216,216,216, 0.5)" : "white")};
  font-size: 2rem;
  background: transparent;
  color: ${props => (props.isDisabled ? "rgba(216,216,216, 0.5)" : "white")};
  transition: 0.15s ease transform;
  margin-top: 2rem;
  max-width: 13rem;
  align-self: center;

  &:hover {
    ${props =>
      props.isDisabled
        ? ""
        : css`
            transform: translate(0, -0.2rem);
          `};
  }

  &:active {
    transform: translate(0, 0.3rem);
  }
`;

const PortWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  width: 80%;
  align-self: center;
`;

const PortConnectionInput = styled(NumberInput)`
  width: 100%;
`;

const StartConnectionButton = styled(StartButton)`
  margin-left: 1.5rem;
  min-width: 7rem;
`;

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      messages: [],
      isConnected: false,
      port: 4242
    };

    this.connection = null;
  }

  onSendPing() {
    this.setState(prev => {
      return [...prev.messages, { msg: "Sent ping!", isError: false }];
    });
  }

  onReceivePong(pong) {
    this.setState(prev => {
      return [...prev.messages, { msg: `Received ${pong}!`, isError: false }];
    });
  }

  onError(error) {
    this.setState(prev => {
      return [...prev.messages, { msg: `Received ${error}!`, isError: true }];
    });
  }

  onSocketClose(code) {
    this.setState(prev => {
      return [...prev.messages, { msg: `Socket closed with code ${code}!`, isError: true }];
    });
  }

  handlePortChange(value) {
    this.setState({ port: value });
  }

  handleConnectionChange() {
    if (this.state.isConnected) {
      this.connection.close();
    } else {
      this.connection = new WebSocket(`ws://localhost:${this.state.port}/commands`);
      if (this.connection.readyState !== 1) {
        this.onError("Websocket connection failed");
        return;
      }

      this.connection.onclose = event => this.onSocketClose(event.code);
      this.connection.onmessage = event => this.onReceivePong(event.data);
      this.connection.onerror = event => this.onError(event);
    }

    this.setState(prev => {
      return { isConnected: !prev.isConnected };
    });
  }

  render() {
    const { port, isConnected } = this.state;
    return (
      <Header>
        <div className="container">
          <div className="row center-xs">
            <div className="col-xs-12">
              <Content>
                <Title className="App-title">Websocket example</Title>
                <PortWrapper>
                  <PortConnectionInput onChange={value => this.handlePortChange(value)} value={port} isEnabled={!isConnected} />
                  <StartConnectionButton onClick={() => this.handleConnectionChange()} />
                </PortWrapper>
                <MessagesList>
                  {this.state.messages.map((message, index) => {
                    return message.isError ? <InfoMessage key={index}>↳ {message.msg}</InfoMessage> : <ErrorMessage key={index}>↳ {message.msg}</ErrorMessage>;
                  })}
                </MessagesList>
                <PingButton disabled={!this.state.isConnected} isDisabled={!this.state.isConnected}>
                  send ping
                </PingButton>
              </Content>
            </div>
          </div>
        </div>
      </Header>
    );
  }
}

export default App;
