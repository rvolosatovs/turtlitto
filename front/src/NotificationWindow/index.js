import React from "react";
import styled, { css } from "styled-components";
import faTimes from "@fortawesome/fontawesome-free-solid/faTimes";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import PropTypes from "prop-types";
import notificationTypes from "./notificationTypes";

const Window = styled.div`
  ${props => props.background};
  border-style: solid;
  border-width: 0.4rem;
  width: 90%;
  position: relative;
  margin: 0 auto;
`;

const ToolBar = styled.div`
  display: flex;
  padding: 1.5rem 1rem;
  align-items: center;
`;

const Title = styled.h3`
  text-align: center;
  font-size: 4rem;
  text-transform: uppercase;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex-grow: 1;
  margin: 0;
`;

const Text = styled.p`
  padding: 0 1rem;
  text-align: center;
  font-size: 3rem;
  overflow-wrap: break-word;
`;

const CloseButton = styled.button`
  border: none;
  background-color: transparent;
  font-size: 2rem;
`;

const getBackground = type => {
  switch (type) {
    case notificationTypes.SUCCESS:
      return css`
        background: ${props => props.theme.notificationSuccess};
      `;
    case notificationTypes.WARNING:
      return css`
        background: ${props => props.theme.notificationWarning};
      `;
    case notificationTypes.ERROR:
      return css`
        background: ${props => props.theme.notificationError};
      `;
    default:
      throw new Error("Unknown notification type");
  }
};

/**
 * Creates a notification window
 * Author: S.A. Tanja
 * Author: G.W. van der Heijden
 * Author: B. Afonins
 *
 * Props:
 * - notificationType: the title and color selection
 * - onDismiss: pass a function that causes the component to no longer be rendered
 * - message: the message of the notification window
 * - children: possible extra elements that could be included (buttons, etc)
 */
const NotificationWindow = props => {
  const background = getBackground(props.notificationType);

  return (
    <Window background={background}>
      <ToolBar>
        <Title>{props.notificationType}</Title>
        <CloseButton onClick={() => props.onDismiss()}>
          <FontAwesomeIcon icon={faTimes} color="black" size="4x" />
        </CloseButton>
      </ToolBar>
      <Text>{props.message}</Text>
      <div>{props.children}</div>
    </Window>
  );
};

NotificationWindow.propTypes = {
  notificationType: PropTypes.oneOf(Object.values(notificationTypes))
    .isRequired,
  onDismiss: PropTypes.func.isRequired,
  message: PropTypes.string.isRequired
};

export default NotificationWindow;
