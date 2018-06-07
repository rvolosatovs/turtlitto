import React, { Fragment } from "react";
import styled, { css } from "styled-components";
import PropTypes from "prop-types";
import connectionTypes from "../connectionTypes";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import faSpinner from "@fortawesome/fontawesome-free-solid/faSpinner";

const Bar = styled.div`
  ${props => props.background};
  padding: 0.6rem 3rem 0.6rem 0.6rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  position: relative;
`;

const Message = styled.p`
  text-transform: uppercase;
  margin: 0;
  display: block;
  font-size: 1.4rem;
`;

const IconWrapper = styled.div`
  padding-left: 2rem;
  position: absolute;
  top: 50%;
  right: 1rem;
  transform: translateY(-50%);
`;

const getBackground = type => {
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

const getContent = type => {
  switch (type) {
    case connectionTypes.CONNECTING:
      return (
        <Fragment>
          <Message>connecting...</Message>
          <IconWrapper>
            <FontAwesomeIcon size="2x" spin icon={faSpinner} />
          </IconWrapper>
        </Fragment>
      );
    case connectionTypes.CONNECTED:
      return <Message>connected</Message>;
    case connectionTypes.DISCONNECTED:
      return <Message>disconnected</Message>;
    default:
      throw new Error("Unknown connection status type");
  }
};

const ConnectionBar = props => {
  const background = getBackground(props.connectionStatus);
  const content = getContent(props.connectionStatus);
  return (
    <Bar className={props.className} background={background}>
      {content}
    </Bar>
  );
};

ConnectionBar.propTypes = {
  connectionStatus: PropTypes.oneOf(Object.values(connectionTypes)).isRequired
};

export default ConnectionBar;
