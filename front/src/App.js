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
import AuthenticationScreen from "./AuthenticationScreen";

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
      token: "dummy",
      turtles: {},
      notifications: [
        { notificationType: "error", message: "Pants on fire" },
        { notificationType: "success", message: "Rendering Notifications" }
      ],
      loggedIn: false
    };
    this.connection = null;
  }

  connect() {
    const l = window.location;
    this.connection = new WebSocket(
      `${l.protocol === "https:" ? "wss" : "ws"}://user:${this.state.token}@${
        l.host
      }/api/v1/state`
    );
    this.connection.onclose = event => this.onConnectionClose(event);
    this.connection.onerror = event => this.onConnectionError(event);
    this.connection.onmessage = event => this.onConnectionMessage(event);
    this.connection.onopen = event => this.onConnectionOpen(event);

    this.setState({ connectionStatus: connectionTypes.CONNECTING });
  }

  componentDidMount() {
    this.connect();
  }

  componentWillUnmount() {
    this.connection.close();
    if (this.timer !== null) {
      clearTimeout(this.timer);
    }
  }

  onConnectionClose(event) {
    this.setState({ connectionStatus: connectionTypes.DISCONNECTED });
    this.timer = setTimeout(() => {
      this.timer = null;
      this.connect();
    }, 1000);
  }

  onConnectionError(event) {
    this.setState({ connectionStatus: connectionTypes.DISCONNECTED });
  }

  onConnectionMessage(event) {
    const data = JSON.parse(event.data);
    this.setState(prev => {
      const turtleChanges = Object.keys(data.turtles).reduce((acc, id) => {
        if (prev.turtles[id] === undefined) {
          data.turtles[id].enabled = false;
          acc[id] = { $set: data.turtles[id] };
        } else {
          acc[id] = { $merge: data.turtles[id] };
        }
        return acc;
      }, {});
      const turtles = update(prev.turtles, turtleChanges);
      return { turtles };
    });
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

  /**
   * The function AuthenticationScreen will call when a token has been
   * submitted.
   * @param token The token received from AuthenticationScreen.
   * @param onIncorrectToken Callback from AuthenticationScreen to update the
   * AuthenticationScreen in case the token was incorrect.
   */
  onSubmit(token, onIncorrectToken) {
    //TODO: Implement a correct version
    if (token === "techunited") {
      this.setState({ loggedIn: true, token: token });
    } else {
      onIncorrectToken();
    }
  }

  getNextNotification() {
    const { notifications } = this.state;
    return notifications.length > 0 ? notifications[0] : null;
  }

  render() {
    const { activePage, turtles, loggedIn, connectionStatus } = this.state;

    return (
      <ThemeProvider theme={theme}>
        {loggedIn ? (
          <Container>
            {activePage === "refbox" && (
              <div>
                <RefboxField isPenalty={false} token={this.state.token} />
                <RefboxSettings token={this.state.token} />
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
                  <Settings turtles={turtles} token={this.state.token} />
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
              token={this.state.token}
            />
          </Container>
        ) : (
          <AuthenticationScreen
            onSubmit={(token, callback) => this.onSubmit(token, callback)}
          />
        )}
      </ThemeProvider>
    );
  }
}

export default App;
