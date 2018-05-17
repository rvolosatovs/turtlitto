import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import styled from "styled-components";
import SRRButton from "./SRRButton";
import Turtle from "./Turtle";
import TurtleEnableBar from "./TurtleEnableBar";
import NotificationWindow from "./NotificationWindow";
import RefboxField from "./RefboxField";

const AppWrap = styled.div`
  height: 100vh;
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

const Middle = styled.div`
  position: relative
  bottom: 90%
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

const composeSendCommand = payload => {
  return JSON.stringify(payload);
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
      host: "localhost",
      activePage: "settings", // Acceptable values: settings, refbox
      turtles: [
        {
          battery: 100,
          enabled: false,
          home: "Yellow home",
          id: 1,
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          battery: 100,
          enabled: false,
          home: "Yellow home",
          id: 2,
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          battery: 100,
          enabled: false,
          home: "Yellow home",
          id: 3,
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          battery: 100,
          enabled: false,
          home: "Yellow home",
          id: 4,
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          battery: 100,
          enabled: false,
          home: "Yellow home",
          id: 5,
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          battery: 100,
          enabled: false,
          home: "Yellow home",
          id: 6,
          role: "INACTIVE",
          team: "Magenta"
        }
      ]
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
    const toSend = composeSendCommand(command);
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

  disableTurtle(position) {
    this.setState((prevState, props) => {
      const turtles = prevState.turtles;
      turtles[position]["enabled"] = false;
      return { turtles: turtles };
    });
  }

  enableTurtle(position) {
    console.log(position);
    this.setState((prevState, props) => {
      const turtles = prevState.turtles;
      turtles[position]["enabled"] = true;
      return { turtles: turtles };
    });
  }

  render() {
    return (
      <AppWrap id="AppWrap">
        {this.state.activePage === "settings" && (
          <div>
            <TurtleBar
              turtles={this.state.turtles}
              onEnable={position => {
                this.enableTurtle(position);
              }}
              onDisable={position => {
                this.disableTurtle(position);
              }}
            />
            <NotificationWindow
              id="NotificationWindow"
              backgroundColor="Tomato"
              NotificationType="Critical Error"
            >
              Turtle 2 died
            </NotificationWindow>
            <Middle id="Middle">
              {this.state.turtles
                .filter(turtle => {
                  return turtle.enabled;
                })
                .map(turtle => {
                  return <Turtle turtle={turtle} />;
                })}
              {showConnected(this.state.isConnected)}
            </Middle>
            <SRRButton
              buttonText={getButtonText(this.state.isConnected)}
              onClick={() => {
                this.handleConnectionStatusChange("mount");
              }}
              enabled={true}
            />
          </div>
        )}
        {this.state.activePage === "refbox" && <RefboxField />}
        <Footer id="Footer">
          <StartButton
            buttonText={<AugmentedText>&#9658;</AugmentedText>}
            onClick={() => {
              this.onSend("start");
            }}
            enabled={true}
          />
          {this.state.activePage === "refbox" && (
            <SettingsButton
              buttonText={<AugmentedText>Settings</AugmentedText>}
              onClick={() => {
                this.setState({ activePage: "settings" });
              }}
              enabled={true}
            />
          )}
          {this.state.activePage === "settings" && (
            <SettingsButton
              buttonText={<AugmentedText>Refbox</AugmentedText>}
              onClick={() => {
                this.setState({ activePage: "refbox" });
              }}
              enabled={true}
            />
          )}
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
