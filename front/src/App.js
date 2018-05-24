import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import { ThemeProvider } from "styled-components";
import theme from "./theme";

import BottomBar from "./BottomBar";
import NotificationWindow from "./NotificationWindow";
import RefboxField from "./RefboxField";
import RefboxSettings from "./RefboxSettings";
import Settings from "./Settings";

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      activePage: "refbox",
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
          enabled: true,
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
          enabled: true,
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

  render() {
    const { activePage, turtles, connectionStatus, notifications } = this.state;
    return (
      <ThemeProvider theme={theme}>
        <div style={{ height: "100vh" }}>
          {activePage === "refbox" && (
            <div>
              <RefboxField />
              <RefboxSettings />
            </div>
          )}
          {activePage === "settings" && <Settings turtles={turtles} />}
          {notifications.map(notification => {
            return (
              <NotificationWindow
                {...notification}
                onDismiss={() =>
                  this.setState(oldState => {
                    return {
                      notifications: oldState.notifications.filter(
                        n => n !== notification
                      )
                    };
                  })
                }
              />
            );
          })}
          <BottomBar
            activePage={activePage}
            changeActivePage={page => this.setState({ activePage: page })}
            onSend={message => this.onSend(message)}
            connectionStatus={connectionStatus}
          />
        </div>
      </ThemeProvider>
    );
  }
}

export default App;
