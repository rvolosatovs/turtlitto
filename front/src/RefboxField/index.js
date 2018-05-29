import React from "react";
import RefboxButton from "./RefboxButton";
import styled, { css } from "styled-components";
import PropTypes from "prop-types";

/**
 * Constructs a refbox for a team consisting of 6 buttons:
 * - Kick off (KO)
 * - Free kick (FK)
 * - Goal kick (GK)
 * - Throw in (TI)
 * - Corner (C)
 * - Penalty (P)
 * Also takes care of handling commands that the buttons represents
 * Author: S.A. Tanja
 * Author: G.W. van der Heijden
 *
 * Props:
 *  - teamColor: the team color (cyan or magenta)
 *  - onClick: a function on what to do when the button is pressed
 *  - isPenalty: a boolean which indicates whether to go into penalty mode
 */
const RefboxField = props => {
  return (
    <Refbox isPenalty={props.isPenalty}>
      {tags(props).map(tag => {
        return (
          <RefboxButton
            key={tag}
            teamColor={props.teamColor}
            onClick={() => {
              props.onClick(tag, props.teamColor);
            }}
          >
            {tag}
          </RefboxButton>
        );
      })}
    </Refbox>
  );
};

/**
 * Places the buttons in 2 columns and rows of 3
 * Allows for multiple refboxes to be next to each other (float : left)
 */
const Refbox = styled.div`
  display: grid;
  ${props =>
    props.isPenalty
      ? css`
          grid-template-columns: repeat(1, 1fr);
        `
      : css`
          grid-template-columns: repeat(2, 1fr);
        `};
  margin: auto;
  width: 16rem;
  float: center;
`;

const tags = props => {
  if (props.isPenalty) {
    return ["Soft", "Medium", "Hard"];
  } else {
    return ["KO", "FK", "GK", "TI", "C", "P"];
  }
};

RefboxField.propType = {
  isPenalty: PropTypes.bool.isRequired,
  teamColor: PropTypes.oneOf(["cyan", "magenta"]).isRequired,
  onClick: PropTypes.func.isRequired
};

export default RefboxField;
