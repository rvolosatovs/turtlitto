import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import { Grid, Row, Col } from "react-flexbox-grid";
import styled, { css, ThemeProvider } from "styled-components";
import format from "date-fns/format";
import theme from "./theme";

import StartButton from "./StartButton";
import InputField from "./InputField";
import Dropdown from "./Dropdown";
import SRRButton from "./SRRButton";
import SRRSwitch from "./SRRSwitch";

const Header = styled.header`
  height: 100vh;
  background: ${props => props.theme.background};
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
  color: ${props => props.theme.title};
`;

const MessagesList = styled.ul`
  list-style-type: none;
  height: 40rem; // dont do that
  margin: 0 auto;
  border: 0.2rem solid ${props => props.theme.secondary};
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

const InfoMessage = Message.extend`
  color: ${props => props.theme.info};
`;

const ErrorMessage = Message.extend`
  color: ${props => props.theme.error};
`;

const SuccessMessage = Message.extend`
  color: ${props => props.theme.success};
`;

const PingButton = styled.button`
  padding: 0.5rem 1rem;
  border-radius: 1.5rem;
  text-transform: uppercase;
  border: 0.2rem solid
    ${props =>
      props.isDisabled
        ? props.theme.baseButtonDisabled
        : props.theme.baseButton};
  font-size: 2rem;
  background: transparent;
  color: ${props =>
    props.isDisabled ? props.theme.baseButtonDisabled : props.theme.baseButton};
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

const ConnectionWrapper = styled.div`
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  width: 80%;
  align-self: center;
`;

const PortConnectionInput = styled(InputField)`
  flex-basis: 70%;
`;

const HostConnectionInput = styled(InputField)`
  flex-basis: 70%;
`;

const StartConnectionButton = styled(StartButton)`
  flex-basis: 25%;
`;

const MessageTypes = Object.freeze({ INFO: 1, ERROR: 2, SUCCESS: 3 });

const styleOutputMessage = (message, key) => {
  switch (message.type) {
    case MessageTypes.INFO:
      return <InfoMessage key={key}>{message.msg}</InfoMessage>;
    case MessageTypes.ERROR:
      return <ErrorMessage key={key}>{message.msg}</ErrorMessage>;
    case MessageTypes.SUCCESS:
      return <SuccessMessage key={key}>{message.msg}</SuccessMessage>;
    default:
      throw new Error("Unknown message type");
  }
};

class AppSocketExample extends Component {
  constructor(props) {
    super(props);

    this.state = {
      messages: [],
      isConnected: false,
      port: "4242",
      host: "localhost"
    };

    this.connection = null;
  }

  onSendPing() {
    this.connection.send("ping");
    this.setState(prev => {
      return {
        messages: [
          ...prev.messages,
          {
            msg: `${format(new Date(), "HH:mm:ss:SS")} > Sent ping!`,
            type: MessageTypes.INFO
          }
        ]
      };
    });
  }

  onReceivePong(pong) {
    this.setState(prev => {
      return {
        messages: [
          ...prev.messages,
          {
            msg: `${format(new Date(), "HH:mm:ss:SS")} > Received ${pong}!`,
            type: MessageTypes.INFO
          }
        ]
      };
    });
  }

  onError(event) {
    const errorMessage = `Cant establish connection to ws://${
      this.state.host
    }:${this.state.port}/commands`;

    this.setState(prev => {
      return {
        messages: [
          ...prev.messages,
          {
            msg: `${format(
              new Date(),
              "HH:mm:ss:SS"
            )} > Error: ${errorMessage}`,
            type: MessageTypes.ERROR
          }
        ]
      };
    });
  }

  onSocketClose(code) {
    if (this.state.isConnected) {
      this.setState(prev => {
        return {
          messages: [
            ...prev.messages,
            {
              msg: `${format(
                new Date(),
                "HH:mm:ss:SS"
              )} > Socket closed with code ${code}!`,
              type: MessageTypes.ERROR
            }
          ],
          isConnected: !prev.isConnected
        };
      });
    }
  }

  onSocketOpen() {
    if (!this.state.isConnected) {
      this.setState(prev => {
        return {
          messages: [
            {
              msg: `${format(
                new Date(),
                "HH:mm:ss:SS"
              )} > Connection established!`,
              type: MessageTypes.SUCCESS
            }
          ],
          isConnected: !prev.isConnected
        };
      });
    }
  }

  handleConnectionChange(value, name) {
    this.setState({ [name]: value });
  }

  handleConnectionStatusChange() {
    if (this.state.isConnected) {
      this.connection.close();
    } else {
      try {
        this.connection = new WebSocket(
          `ws://${this.state.host}:${this.state.port}/commands`
        );
        this.connection.onclose = event => this.onSocketClose(event.code);
        this.connection.onmessage = event => this.onReceivePong(event.data);
        this.connection.onerror = event => this.onError(event);
        this.connection.onopen = event => this.onSocketOpen();
      } catch (e) {
        this.onError(e.message);
      }
    }
  }

  render() {
    const { port, isConnected, host } = this.state;
    return (
      <ThemeProvider theme={theme}>
        <Header>
          <Grid>
            <Row center="xs">
              <Col xs={12}>
                <Content>
                  <Title className="App-title">Websocket example</Title>
                  <ConnectionWrapper>
                    <PortConnectionInput
                      onChange={value =>
                        this.handleConnectionChange(value, "port")
                      }
                      value={port}
                      isDisabled={isConnected}
                      placeholder="Port"
                    />
                    <HostConnectionInput
                      onChange={value =>
                        this.handleConnectionChange(value, "host")
                      }
                      value={host}
                      isDisabled={isConnected}
                      placeholder="Host"
                    />
                    <StartConnectionButton
                      onClick={() => this.handleConnectionStatusChange()}
                      isRunning={isConnected}
                    />
                  </ConnectionWrapper>
                  <MessagesList>
                    {this.state.messages.map((message, index) => {
                      return styleOutputMessage(message, index);
                    })}
                  </MessagesList>
                  <PingButton
                    onClick={() => this.onSendPing()}
                    disabled={!this.state.isConnected}
                    isDisabled={!this.state.isConnected}
                  >
                    send ping
                  </PingButton>
                  <Dropdown
                    values={["Henk", "Wessel"]}
                    currentValue={"Wessel"}
                    onChange={value => {
                      console.log(value);
                    }}
                    enabled={true}
                  />
                </Content>
                <SRRButton
                  buttonText={"hurrdurr Ahma Button"}
                  onClick={() => {
                    console.log("harrrr");
                  }}
                  enabled={true}
                />
                <SRRSwitch
                  currentValue={true}
                  buttonText={"hurrdurr Ahma Switch"}
                  onChange={() => {
                    console.log("hurrrr");
                    this.currentValue = !this.currentValue;
                  }}
                  enabled={true}
                />
              </Col>
            </Row>
          </Grid>
        </Header>
      </ThemeProvider>
    );
  }
}

export default AppSocketExample;
