import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import styled from "styled-components";
import SRRButton from "./SRRButton";
import Turtle from "./Turtle";
import TurtleEnableBar from "./TurtleEnableBar";
import NotificationWindow from "./NotificationWindow";

const AppWrap = styled.div`
  width = 100%;
  height = 100%;
`;

const TurtleBar = styled(TurtleEnableBar)`
  display: flex;
  width: 100%;
  justify-content: space-between;
  overflow-x: auto;
  scrollbar: hidden;
  align-content: space-between;
  border-style: solid;
  border-width: 0px 0px 2px 0px;
  margin-bottom: 2px;
  background-color: black;
`;

const Footer = styled.footer`
  position: fixed;
  background-color: white;
  bottom: 0;
  width: 100%;
  min-height: 10%;
  border-style: solid;
  border-width: 2px 0px 0px 0px;
  display: flex;
  justify-content: space-between;
  margin: 0px;
`;

const StartButton = styled(SRRButton)`
  width: 32%;
  font-size: 5vmin;
  flex: 1;
`;

const SettingsButton = styled(SRRButton)`
  width: 32%;
  font-size: 5vmin;
  flex: 1;
`;

const StopButton = styled(SRRButton)`
  width: 32%;
  font-size: 5vmin;
  flex: 1;
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
    if (
      this.state.isConnected ||
      (status === "unmount" && this.state.isConnected)
    ) {
      this.connection.close();
    } else if (status !== "unmount" && !this.state.isConnected) {
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
        <TurtleBar />
        <NotificationWindow
          backgroundColor="Tomato"
          NotificationType="Critical Error"
        >
          Turtle 2 died
        </NotificationWindow>
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
