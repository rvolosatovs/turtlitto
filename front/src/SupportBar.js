import React, { Component } from "react";
import styled, { css } from "styled-components";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import faTimes from "@fortawesome/fontawesome-free-solid/faTimes";

/**
 * Hack to show the complete bottom bar on some broken browsers, as is some browsers, the refbox switch button is obscured by the address bar, and scrolling is blocked, so the address bar is never hidden.
 * This is needed for Chrome on Android and Safari on iOS. An close button is added to make it possible to close this bar if it becomes visible by accidental correct standards implementation. If at some point in the future this bar can be removed, rejoice.
 *
 * Author: H.E. van der Laan (unfortunately)
 */
export default class SupportBar extends Component {
  constructor(props) {
    super(props);
    /* The following browsers need a support bar:
     * - Safari on iOS
     * - Chrome on Android
     * See https://developer.chrome.com/multidevice/user-agent
     */
    const isAndroidChrome = /Android .* Chrome\/[.0-9]* Mobile/.test(
      navigator.userAgent
    );
    const isIosSafari = /iP.* Version/.test(navigator.userAgent);
    const isIpad = /iPad.* Version/.test(navigator.userAgent);
    this.state = { show: isAndroidChrome || isIosSafari, isIpad: isIpad };
    this.onResize = () => this.forceUpdate();
    window.addEventListener("resize", this.onResize);
  }

  componentWillUnmount() {
    window.removeEventListener("resize", this.onResize);
  }

  render() {
    // Don't show the SupportBar in landscape mode
    if (this.state.show) {
      return (
        <Bar isIpad={this.state.isIpad}>
          <span>Soccer Robot Remote</span>
          <CloseButton onClick={() => this.setState({ show: false })}>
            <FontAwesomeIcon icon={faTimes} />
          </CloseButton>
        </Bar>
      );
    } else {
      return null;
    }
  }
}

const Bar = styled.div`
  align-items: center;
  display: flex;
  font-size: 3rem;
  ${props =>
    props.isIpad
      ? css`
          height: 4vh;
        `
      : css`
          height: 9vh;
        `};
  justify-content: space-around;
`;

const CloseButton = styled.button`
  border: none;
`;
