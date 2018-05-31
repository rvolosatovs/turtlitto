import React, { Component, Fragment } from "react";
import "./App.css";
import "normalize.css";
import styled, { ThemeProvider } from "styled-components";
import theme from "./theme";

import BottomBar from "./BottomBar";
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
      connectionStatus: "connected",
      turtles: [
        {
          id: 1,
          enabled: false,
          battery: 66,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          id: 2,
          enabled: false,
          battery: 42,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          id: 3,
          enabled: false,
          battery: 42,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          id: 4,
          enabled: false,
          battery: 100,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          id: 5,
          enabled: false,
          battery: 4,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        },
        {
          id: 6,
          enabled: false,
          battery: 0,
          home: "Yellow home",
          role: "INACTIVE",
          team: "Magenta"
        }
      ],
      notifications: [
        { notificationType: "error", message: "Pants on fire" },
        { notificationType: "success", message: "Rendering Notifications" }
      ]
    };
  }

  onSend(message) {
    console.log(`Sent: ${message}`);
  }

  onTurtleEnableChange(position) {
    this.setState(prev => {
      const turtles = prev.turtles.map((turtle, index) => {
        if (index === position) {
          return {
            ...turtle,
            enabled: !turtle.enabled
          };
        }

        return turtle;
      });

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
              <RefboxField />
              <RefboxSettings />
            </div>
          )}
          {activePage === "settings" && (
            <Fragment>
              <TurtleEnableBar
                turtles={turtles.map(turtle => {
                  return {
                    id: turtle.id,
                    enabled: turtle.enabled
                  };
                })}
                onTurtleEnableChange={position =>
                  this.onTurtleEnableChange(position)
                }
              />
              <ScrollableContent>
                <Settings turtles={turtles.filter(turtle => turtle.enabled)} />
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
            onSend={message => this.onSend(message)}
            connectionStatus={connectionStatus}
          />
        </Container>
      </ThemeProvider>
    );
  }
}

export default App;
