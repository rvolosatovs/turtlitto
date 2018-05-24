import styled from "styled-components";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import faPlay from "@fortawesome/fontawesome-free-solid/faPlay";
import faStop from "@fortawesome/fontawesome-free-solid/faStop";
import faCog from "@fortawesome/fontawesome-free-solid/faCog";
import faThLarge from "@fortawesome/fontawesome-free-solid/faThLarge";
import React from "react";
import PropTypes from "prop-types";
import pageTypes from "./pageTypes";

/**
 * All of the props depend on the implementation in App.js. App.js should include WebSockets, with an onSend command,
 * which sends a command,a state variable which has the values "settings" and "refbox",
 * and the implementation of the functions passed to BottomBar.
 * Author: T.T.P. Franken
 * Author: B. Afonins
 *
 * Props:
 * - onSend: function to send a command
 * - changeActivePage: function to change the active page
 * - activePage: a string indicating the current active page
 */

const BottomBar = props => {
  const { changeActivePage, onSend, activePage } = props;
  const isSettingsPage = activePage === pageTypes.SETTINGS;
  return (
    <Bar className={props.className}>
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
      <Button
        id="bottom-bar__stop-button"
        onClick={() => onSend("stop")}
        enabled
      >
        <FontAwesomeIcon icon={faStop} color="red" />
      </Button>
    </Bar>
  );
};

BottomBar.propTypes = {
  onSend: PropTypes.func.isRequired,
  changeActivePage: PropTypes.func.isRequired,
  activePage: PropTypes.oneOf(Object.values(pageTypes))
};

const Bar = styled.div`
  position: fixed;
  background-color: white;
  bottom: 0;
  width: 100%;
  height: 10%;
  border-style: solid;
  border-width: 0.2rem 0rem 0rem 0rem;
  display: flex;
  justify-content: space-between;
  margin: 0rem;
`;

const Button = styled.button`
  background: ${props => props.theme.bottomBarButton};
  width: 32%;
  font-size: 5rem;
  flex: 1;

  &:active {
    background: ${props => props.theme.bottomBarButtonActive};
  }
`;

export default BottomBar;
