import PropTypes from "prop-types";
import React from "react";
import styled from "styled-components";

/**
 * A button part of the refbox. Does a callback to its parent to handle onClick events
 * Props:
 *  - children: the text the button needs to display
 *  - color: Contains the background color of the button
 *
 * Author: S.A. Tanja
 */
const RefboxButton = props => {
  return (
    <Button
      color={props.color}
      onClick={() => {
        props.onClick(props.color, props.children);
      }}
    >
      {props.children}
    </Button>
  );
};

RefboxButton.propTypes = {
  color: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired
};

/**
 * A styled component defining the style of the refbox buttons
 */
const Button = styled.button`
  background-color: ${props => props.color};
  color: black;
  width: 206px;
  height: 206px;
  border: solid;
  border-width: 4px;
  text-align: center;
  font-size: 100px;
`;

export default RefboxButton;
