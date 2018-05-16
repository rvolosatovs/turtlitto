import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import styled from "styled-components";
import SRRButton from "./SRRButton";
import Turtle from "./Turtle";
import format from "date-fns/format";

const AppWrap = styled.div`\
  width = 100%;
  height = 100%;
`;

const Footer = styled.footer`
  position: fixed;
  bottom: 0;
  width: 100%;
  min-width: 100px;
  height: 10%;
  border-style: solid;
  border-width: 5px 0px 0px 0px;
  display: flex;
  justify-content: space-between;
  margin: 0px;
`;

const StartButton = styled(SRRButton)`
  width: 32%;
  font-size: 4vw;
  content: "\&#9658;";
`;

const SettingsButton = styled(SRRButton)`
  width: 32%;
  font-size: 3vw;
`;

const StopButton = styled(SRRButton)`
  width: 32%;
  font-size: 4vw;
`;

const AugmentedText = styled.p`
  margin: 0px;
  padding: 0px;
`;

const composeSendCommand = (type, payload) => {
  const command = {
    type: type,
    message_id: "undefined",
    payload: payload
  };
  return JSON.stringify(command);
};

const showConnected = status => {
  if (status === true) {
    return <p>Connection: connected</p>;
  }
  return <p>Connection: disconnected</p>;
};

const getButtonText = status => {
  if (status === true) {
    return <p>Disconnect</p>;
  }
  return <p>Reconnect</p>;
};

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      //messages: [],
      isConnected: false,
      port: "4242",
      host: "localhost"
    };

    this.connection = null;
  }

  componentDidMount() {
    this.handleConnectionStatusChange("mount");
  }

  componentWillUnmount() {
    this.handleConnectionStatusChange("unmount");
  }

  handleConnectionStatusChange(status) {
    if (this.state.isConnected || status === "unmount") {
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

  onError(event) {
    const errorMessage = `Cant establish connection to ws://${
      this.state.host
    }:${this.state.port}/commands`;
    console.log(errorMessage);
  }

  onSend(command) {
    const toSend = composeSendCommand("command", command);
    if (this.state.isConnected) {
      this.connection.send(toSend);
      console.log("Sending is a success!");
      console.log(toSend);
    } else {
      console.log("there is no connection at the moment.");
    }
  }

  onReceivePong(pong) {
    console.log(pong);
  }

  onSocketClose(code) {
    if (this.state.isConnected) {
      this.setState(prev => {
        return {
          isConnected: !prev.isConnected
        };
      });
    }
    this.forceUpdate();
  }

  onSocketOpen() {
    if (!this.state.isConnected) {
      this.setState(prev => {
        return {
          isConnected: !prev.isConnected
        };
      });
    }
    this.forceUpdate();
  }

  render() {
    return (
      <AppWrap>
        <Turtle
          turtle={{
            battery: 100,
            editable: true,
            role: "Goalkeeper",
            home: "Yellow home",
            team: "Magenta",
            id: 2
          }}
        />
        {showConnected(this.state.isConnected)}
        <SRRButton
          buttonText={getButtonText(this.state.isConnected)}
          onClick={() => {
            this.handleConnectionStatusChange("mount");
          }}
          enabled={true}
        />
        <Footer>
          <StartButton
            buttonText={<AugmentedText>&#9658;</AugmentedText>}
            onClick={() => {
              this.onSend("start");
            }}
            enabled={true}
          />
          <SettingsButton
            buttonText={<AugmentedText>Settings</AugmentedText>}
            onClick={() => {}}
            enabled={true}
          />
          <StopButton
            buttonText={<AugmentedText>&#9724;</AugmentedText>}
            onClick={() => {
              this.onSend("stop");
            }}
            enabled={true}
          />
        </Footer>
      </AppWrap>
    );
  }
}

export default App;
