import styled, { css } from "styled-components";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import faPlay from "@fortawesome/fontawesome-free-solid/faPlay";
import faStop from "@fortawesome/fontawesome-free-solid/faStop";
import faCog from "@fortawesome/fontawesome-free-solid/faCog";
import faThLarge from "@fortawesome/fontawesome-free-solid/faThLarge";
import React from "react";
import PropTypes from "prop-types";
import pageTypes from "./pageTypes";
import connectionTypes from "./connectionTypes";
import sendToServer from "../sendToServer";

const getConnectionStatusBackground = type => {
  switch (type) {
    case connectionTypes.CONNECTING:
      return css`
        background: ${props => props.theme.notificationWarning};
      `;
    case connectionTypes.CONNECTED:
      return css`
        background: ${props => props.theme.notificationSuccess};
      `;
    case connectionTypes.DISCONNECTED:
      return css`
        background: ${props => props.theme.notificationError};
      `;
    default:
      throw new Error("Unknown connection status type");
  }
};

/**
 * All of the props depend on the implementation in App.js. App.js should include WebSockets,
 * a state variable which has the values "settings" and "refbox",
 * and the implementation of the functions passed to BottomBar.
 * Author: T.T.P. Franken
 * Author: B. Afonins
 * Author: G.M. van der Sanden
 *
 * Props:
 * - changeActivePage: function to change the active page
 * - activePage: a string indicating the current active page
 * - connectionStatus: a boolean indicating whether the client is connected to the TRC
 * - token: a string which holds the password needed to connect to the SRRS
 */

const BottomBar = props => {
  const { changeActivePage, activePage, connectionStatus, token } = props;
  const isSettingsPage = activePage === pageTypes.SETTINGS;
  const connectionStatusBackground = getConnectionStatusBackground(
    connectionStatus
  );

  return (
    <Bar className={props.className}>
      <ConnectionBar background={connectionStatusBackground}>
        <ConnectionStatus>{connectionStatus}</ConnectionStatus>
      </ConnectionBar>
      <ButtonsWrapper>
        <ButtonColumn>
          <Button
            id="bottom-bar__start-button"
            onClick={() => {
              sendToServer("start", "command", token);
            }}
            enabled
          >
            <FontAwesomeIcon icon={faPlay} />
          </Button>
          <Button
            id="bottom-bar__change-page-button"
            onClick={() =>
              changeActivePage(
                isSettingsPage ? pageTypes.REFBOX : pageTypes.SETTINGS
              )
            }
            enabled
          >
            {isSettingsPage ? (
              <FontAwesomeIcon icon={faThLarge} />
            ) : (
              <FontAwesomeIcon icon={faCog} />
            )}
          </Button>
        </ButtonColumn>
        <StopButton
          id="bottom-bar__stop-button"
          onClick={() => {
            sendToServer("stop", "command", token);
          }}
          enabled
        >
          <FontAwesomeIcon icon={faStop} color="red" />
        </StopButton>
      </ButtonsWrapper>
    </Bar>
  );
};

BottomBar.propTypes = {
  changeActivePage: PropTypes.func.isRequired,
  activePage: PropTypes.oneOf(Object.values(pageTypes)),
  connectionStatus: PropTypes.oneOf(Object.values(connectionTypes)).isRequired
};

const ButtonsWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  margin: 0rem;
  flex: 1;
`;

const ButtonColumn = styled.div`
  width: 50%;
  height: 100%;
  justify-content: space-between;
  margin: 0rem;
  flex: 1;
`;

const Button = styled.button`
  background: ${props => props.theme.bottomBarButton};
  width: 100%;
  height: 50%;
  font-size: 200%;
  flex: 1;

  &:active {
    background: ${props => props.theme.bottomBarButtonActive};
  }
`;

const StopButton = styled.button`
  background: ${props => props.theme.bottomBarButton};
  font-size: 400%;
  flex: 1;

  &:active {
    background: ${props => props.theme.bottomBarButtonActive};
  }
`;

const ConnectionBar = styled.div`
  ${props => props.background};
  color: white;
  margin: 0px;
  text-align: center;
  font-size: 1.2rem;
  padding: 0.6rem;
  text-transform: uppercase;
`;

const Bar = styled.div`
  position: sticky;
  background-color: white;
  top: 100%;
  width: 100%;
  height: 20%;
  display: flex;
  flex-direction: column;
  margin: 0px;
`;

const ConnectionStatus = styled.p`
  margin: 0px;
`;

export default BottomBar;
