import React from "react";
import styled from "styled-components";
import faTimes from "@fortawesome/fontawesome-free-solid/faTimes";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import theme from "../theme";
import PropTypes from "prop-types";

const Window = styled.div`
  background-color: ${props => props.backgroundColor};
  border-style: solid;
  border-width: 0.4rem;
  width: 90%;
  position: relative;
  margin: 0 auto;
`;

const Header = styled.header`
  text-align: center;
  font-size: 5rem;
`;

const Text = styled.p`
  padding: 0 1rem;
  text-align: center;
  font-size: 4rem;
  overflow-wrap: break-word;
`;

const CloseButton = styled.button`
  border: none;
  background-color: transparent;
  position: absolute;
  top: 0;
  right: 0;
  padding: 1rem 2rem;
  font-size: 2rem;
`;

const NotificationTypes = Object.freeze({
  SUCCESS: "SUCCESS",
  WARNING: "WARNING",
  ERROR: "ERROR"
});

/**
 * Creates a notification window
 * Author: S.A. Tanja
 * Author: G.W. van der Heijden
 *
 * Props:
 * - notificationType: the title and color selection
 * - onDismiss: pass a function that causes the component to no longer be rendered
 * - message: the message of the notification window
 * - children: possible extra elements that could be included (buttons, etc)
 */
const NotificationWindow = props => {
  const backgroundColor = () => {
    switch (props.notificationType) {
      case NotificationTypes.SUCCESS:
        return theme.notificationSuccess;
      case NotificationTypes.WARNING:
        return theme.notificationWarning;
      case NotificationTypes.ERROR:
        return theme.notificationError;
      default:
        return "white";
    }
  };
  return (
    <Window backgroundColor={backgroundColor}>
      <Header>{props.notificationType}</Header>
      <CloseButton onClick={() => props.onDismiss()}>
        <FontAwesomeIcon icon={faTimes} color="black" size="4x" />
      </CloseButton>
      <Text>{props.message}</Text>
      <div>{props.children}</div>
    </Window>
  );
};

NotificationWindow.propTypes = {
  notificationType: PropTypes.oneOf(["SUCCESS", "WARNING", "ERROR"]),
  onDismiss: PropTypes.func.isRequired,
  message: PropTypes.string.isRequired
};

export { NotificationTypes, NotificationWindow };
