import React, { Component, Fragment } from "react";
import "./App.css";
import "normalize.css";
import styled, { ThemeProvider } from "styled-components";
import theme from "./theme";
import update from "immutability-helper";

import BottomBar from "./BottomBar";
import connectionTypes from "./BottomBar/connectionTypes";
import NotificationWindow from "./NotificationWindow";
import RefboxField from "./RefboxField";
import RefboxSettings from "./RefboxSettings";
import Settings from "./Settings";
import TurtleEnableBar from "./TurtleEnableBar";

const Container = styled.div`
  height: 100%;
  display: flex;
  flex-direction: column;
`;

const ScrollableContent = styled.div`
  flex: 1;
  overflow-y: auto;
`;

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      activePage: "settings",
      connectionStatus: connectionTypes.DISCONNECTED,
      turtles: {
        1: {
          enabled: false,
          batteryvoltage: 66,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        },
        2: {
          enabled: false,
          batteryvoltage: 42,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        },
        3: {
          enabled: false,
          batteryvoltage: 42,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        },
        4: {
          enabled: false,
          batteryvoltage: 100,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        },
        5: {
          enabled: false,
          batteryvoltage: 4,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        },
        6: {
          enabled: false,
          batteryvoltage: 0,
          homegoal: "Yellow home",
          role: "INACTIVE",
          teamcolor: "Magenta"
        }
      },
      notifications: [
        { notificationType: "error", message: "Pants on fire" },
        { notificationType: "success", message: "Rendering Notifications" }
      ]
    };
    this.connection = null;
  }

  componentDidMount() {
    const l = window.location;
    this.connection = new WebSocket(
      `${l.protocol === "https:" ? "wss" : "ws"}://${l.host}/api/v1/state`
    );
    this.connection.onclose = event => this.onConnectionClose(event);
    this.connection.onerror = event => this.onConnectionError(event);
    this.connection.onmessage = event => this.onConnectionMessage(event);
    this.connection.onopen = event => this.onConnectionOpen(event);

    this.setState({ connectionStatus: connectionTypes.CONNECTING });
  }

  componentWillUnmount() {
    this.connection.close();
  }

  onConnectionClose(event) {
    this.setState({ connectionStatus: connectionTypes.DISCONNECTED });
  }

  onConnectionError(event) {
    this.setState({ connectionStatus: connectionTypes.DISCONNECTED });
  }

  onConnectionMessage(event) {
    console.log(event); // TODO: something useful with this message
  }

  onConnectionOpen(event) {
    this.setState({ connectionStatus: connectionTypes.CONNECTED });
  }

  onSend(message) {
    console.log(`Sent: ${message}`);
  }

  onTurtleEnableChange(id) {
    this.setState(prev => {
      const turtles = update(prev.turtles, { [id]: { $toggle: ["enabled"] } });
      return { turtles };
    });
  }

  onNotificationDismiss() {
    this.setState(prev => {
      return {
        notifications: prev.notifications.slice(1)
      };
    });
  }

  getNextNotification() {
    const { notifications } = this.state;
    return notifications.length > 0 ? notifications[0] : null;
  }

  render() {
    const { activePage, turtles, connectionStatus } = this.state;
    return (
      <ThemeProvider theme={theme}>
        <Container>
          {activePage === "refbox" && (
            <div>
              <RefboxField isPenalty={false} />
              <RefboxSettings />
            </div>
          )}
          {activePage === "settings" && (
            <Fragment>
              <TurtleEnableBar
                turtles={Object.keys(turtles).map(id => {
                  return {
                    id: id,
                    enabled: turtles[id].enabled
                  };
                })}
                onTurtleEnableChange={id => this.onTurtleEnableChange(id)}
              />
              <ScrollableContent>
                <Settings turtles={turtles} />
              </ScrollableContent>
            </Fragment>
          )}
          <NotificationWindow
            onDismiss={() => this.onNotificationDismiss()}
            notification={this.getNextNotification()}
          />
          <BottomBar
            activePage={activePage}
            changeActivePage={page => this.setState({ activePage: page })}
            connectionStatus={connectionStatus}
          />
        </Container>
      </ThemeProvider>
    );
  }
}

export default App;
