import React, { Component, Fragment } from "react";
import "./App.css";
import "normalize.css";
import styled, { ThemeProvider } from "styled-components";
import theme from "./theme";
import update from "immutability-helper";
import { Row, Col } from "react-flexbox-grid";
import pageTypes from "./BottomBar/pageTypes";
import { screenSizes } from "./media";

import Bar from "./BottomBar";
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

const BottomBar = styled(Bar)`
  height: 100%;
`;

const StickyBottomContainer = styled.div`
  position: sticky;
  top: 100%;
  width: 100%;
  margin: 0;
`;

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      activePage: pageTypes.SETTINGS,
      connectionStatus: connectionTypes.DISCONNECTED,
      session: "",
      command: "role_assigner_on",
      turtles: {},
      notifications: [
        { notificationType: "error", message: "Pants on fire" },
        { notificationType: "success", message: "Rendering Notifications" }
      ],
      loggedIn: false,
      authNotification: ""
    };
    this.connection = null;
    this.checkWindowWidth = this.checkWindowWidth.bind(this);
  }

  checkWindowWidth() {
    /*
     * Once enter the tablet mode, make sure
     * we transition to the settings page
    */
    if (window.innerWidth >= screenSizes.md) {
      this.setState({ activePage: pageTypes.SETTINGS });
    }
  }

  connect() {
    const l = window.location;
    this.connection = new WebSocket(
      `${l.protocol === "https:" ? "wss" : "ws"}://${l.host}/api/v1/state`
    );
    this.connection.onclose = event => this.onConnectionClose(event);
    this.connection.onerror = event => this.onConnectionError(event);
    this.connection.onmessage = event => this.onConnectionMessage(event);
    this.connection.onopen = event => this.onConnectionOpen(event);
    window.addEventListener("resize", this.checkWindowWidth);

    this.setState({ connectionStatus: connectionTypes.CONNECTING });
  }

  componentWillUnmount() {
    this.connection.close();
    if (this.timer !== null) {
      clearTimeout(this.timer);
    }
    window.removeEventListener("resize", this.checkWindowWidth);
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
    if (data.turtles !== undefined)
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
    if (data.command !== undefined) this.setState({ command: data.command });
  }

  onConnectionOpen(event) {
    this.connection.send(JSON.stringify(this.state.session));
    this.setState({ connectionStatus: connectionTypes.CONNECTED });
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
   */
  authSubmit(token) {
    const l = window.location;
    console.log(`send authorization to ${l.protocol}//${l.host}/api/v1/auth`);
    fetch(`${l.protocol}//${l.host}/api/v1/auth`, {
      method: "GET",
      headers: new Headers({
        Authorization: "Basic " + btoa(`user:${token}`)
      })
    }).then(response => {
      response
        .text()
        .then(result => {
          if (!response.ok) {
            throw new Error(result);
          }
          this.setState({ loggedIn: true, session: result });
          this.connect();
        })
        .catch(error => {
          this.setState({ authNotification: error.message });
        });
    });
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
            {activePage === pageTypes.SETTINGS && (
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
                  <Settings turtles={turtles} session={this.state.session} />
                </ScrollableContent>
              </Fragment>
            )}
            <StickyBottomContainer>
              <Row bottom="xs">
                <Col md={4} className={"hidden-xs hidden-sm"}>
                  <RefboxField
                    isPenalty={this.state.command === "penalty_demo"}
                    session={this.state.session}
                  />
                </Col>
                <Col md={4} className={"hidden-xs hidden-sm"}>
                  <NotificationWindow
                    onDismiss={() => this.onNotificationDismiss()}
                    notification={this.getNextNotification()}
                  />
                </Col>
                <Col xs={12} md={4} className={"hidden-xs hidden-sm"}>
                  <RefboxSettings session={this.state.session} />
                  <BottomBar
                    activePage={activePage}
                    changeActivePage={() => {}}
                    connectionStatus={connectionStatus}
                    session={this.state.session}
                  />
                </Col>
                <Col xs={12} className={"hidden-md hidden-lg hidden-xl"}>
                  {activePage === pageTypes.REFBOX && (
                    <Fragment>
                      <RefboxField
                        isPenalty={this.state.command === "penalty_demo"}
                        session={this.state.session}
                      />
                      <RefboxSettings session={this.state.session} />
                    </Fragment>
                  )}
                  <NotificationWindow
                    onDismiss={() => this.onNotificationDismiss()}
                    notification={this.getNextNotification()}
                  />
                  <BottomBar
                    activePage={activePage}
                    changeActivePage={page =>
                      this.setState({ activePage: page })
                    }
                    connectionStatus={connectionStatus}
                    session={this.state.session}
                  />
                </Col>
              </Row>
            </StickyBottomContainer>
          </Container>
        ) : (
          <AuthenticationScreen
            notification={this.state.authNotification}
            onSubmit={(token, callback) => this.authSubmit(token, callback)}
            connectionStatus={connectionStatus}
          />
        )}
      </ThemeProvider>
    );
  }
}

export default App;
