import React from "react";
import styled from "styled-components";
import faTimes from "@fortawesome/fontawesome-free-solid/faTimes";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";

const Window = styled.div`
  background-color : ${props => props.backgroundColor}
  border-style: solid;
  border-size: 4px;
  width: 90%;
  position: relative;
  margin-left:auto;
  margin-right:auto;
`;

const NotificationHeader = styled.p`
  text-align: center;
  font-size: 5vmin;
`;

const NotificationText = styled.p`
  padding-left: 10px;
  padding-right: 10px;
  text-align: center;
  font-size: 8vmin;
`;

const NotificationCloseButton = styled.button`
  border: none;
  background-color: transparent;
  position: absolute;
  top: 0;
  right: 0;
  width: 4vmin;
  height: 4vmin;
  &:focus {
    outline: 0;
  }
  font-size: 2vmin;
`;

/**
 * Notification window
 * - Allows setting of background color (no default!)
 * - Allows setting of notification header (props.NotificationType)
 * - Text inside the window is done through props.children
 * - TODO: close notifcation window upon pressing the close button. Class?
 * Author: S.A. Tanja
 */
const NotificationWindow = props => {
  return (
    <Window backgroundColor={props.backgroundColor}>
      <NotificationHeader>{props.NotificationType}</NotificationHeader>
      <NotificationCloseButton
        onClick={() => {
          window.alert("pressed close");
        }}
      >
        <FontAwesomeIcon icon={faTimes} color="Black" size="2x" />
      </NotificationCloseButton>
      <NotificationText>{props.children}</NotificationText>
    </Window>
  );
};
export default NotificationWindow;
