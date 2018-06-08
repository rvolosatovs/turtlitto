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
import sendToServer from "../sendToServer";
import media from "../media";
import ConnectionBar from "./ConnectionBar";

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

  return (
    <Bar className={props.className}>
      <ConnectionBar connectionStatus={connectionStatus} />
      <ButtonsWrapper>
        <ButtonColumn>
          <StartButton
            id="bottom-bar__start-button"
            onClick={() => {
              sendToServer("start", "command", token);
            }}
            enabled
          >
            <FontAwesomeIcon icon={faPlay} />
          </StartButton>
          <ChangePageButton
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
          </ChangePageButton>
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
  margin: 0rem;
`;

const ButtonColumn = styled.div`
  margin: 0rem;
  flex-basis: 50%;
`;

const Button = styled.button`
  background: ${props => props.theme.bottomBarButton};
  width: 100%;
  font-size: 2rem;
  flex: 1;

  &:active {
    background: ${props => props.theme.bottomBarButtonActive};
  }
`;

const StopButton = styled(Button)`
  font-size: 5rem;
  padding: 3rem;
  flex-basis: 50%;
`;

const StartButton = styled(Button)`
  height: 50%;

  ${media.md`
    height: 100%;
    padding: 3rem;
    font-size: 5rem;
  `};
`;

const ChangePageButton = styled(Button)`
  height: 50%;

  ${media.md`
    display: none;
  `};
`;

const Bar = styled.div`
  background-color: white;
  display: flex;
  flex-direction: column;
`;

export default BottomBar;
