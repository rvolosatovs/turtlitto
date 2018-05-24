import styled from "styled-components";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import faPlay from "@fortawesome/fontawesome-free-solid/faPlay";
import faStop from "@fortawesome/fontawesome-free-solid/faStop";
import faCog from "@fortawesome/fontawesome-free-solid/faCog";
import faThLarge from "@fortawesome/fontawesome-free-solid/faThLarge";
import React from "react";
import PropTypes from "prop-types";
import pageTypes from "./pageTypes";
import connectionTypes from "./connectionTypes";

/**
 * All of the props depend on the implementation in App.js. App.js should include WebSockets, with an onSend command,
 * which sends a command,a state variable which has the values "settings" and "refbox",
 * and the implementation of the functions passed to BottomBar.
 * Author: T.T.P. Franken
 * Author: B. Afonins
 * Author: G.M. van der Sanden
 *
 * Props:
 * - onSend: function to send a command
 * - changeActivePage: function to change the active page
 * - activePage: a string indicating the current active page
 * - isConnected: a boolean indicating whether the client is connected to the TRC
 */

const BottomBar = props => {
  const { changeActivePage, onSend, activePage, isConnected } = props;
  const isSettingsPage = activePage === pageTypes.SETTINGS;
  const connecting = isConnected === connectionTypes.CONNECTING;
  const connected = isConnected === connectionTypes.CONNECTED;
  return (
    <Bar>
      <ConnectionBar
        id="connectBar"
        connecting={connecting}
        connected={connected}
      >
        <ConnectionStatus>
          {connecting
            ? connectionTypes.CONNECTING
            : connected
              ? connectionTypes.CONNECTED
              : connectionTypes.DISCONNECTED}
        </ConnectionStatus>
      </ConnectionBar>
      <ButtonsWrapper className={props.className}>
        <ButtonColumn className={props.className}>
          <Button
            id="bottom-bar__start-button"
            onClick={() => onSend("start")}
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
          onClick={() => onSend("stop")}
          enabled
        >
          <FontAwesomeIcon icon={faStop} color="red" />
        </StopButton>
      </ButtonsWrapper>
    </Bar>
  );
};

BottomBar.propTypes = {
  onSend: PropTypes.func.isRequired,
  changeActivePage: PropTypes.func.isRequired,
  activePage: PropTypes.oneOf(Object.values(pageTypes)),
  isConnected: PropTypes.oneOf(Object.values(connectionTypes))
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
  background-color: ${props =>
    props.connecting
      ? props.theme.notificationWarning
      : props.connected
        ? props.theme.notificationSuccess
        : props.theme.notificationError};
  color: white;
  margin: 0px;
  text-align: center;
  font-size: 75%;
  padding: 0.2rem;
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
