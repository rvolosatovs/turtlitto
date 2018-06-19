import React from "react";
import Button from "./RefboxButton";
import styled, { css } from "styled-components";
import PropTypes from "prop-types";
import sendToServer from "../sendToServer";
import media from "../media";

const TAG_VALUES = {
  KO: "kick_off",
  FK: "free_kick",
  GK: "goal_kick",
  TI: "throw_in",
  C: "corner",
  P: "penalty",
  Soft: "soft",
  Medium: "medium",
  Hard: "hard"
};

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
 *  - isPenalty: a boolean which indicates whether to go into penalty mode
 *  - session: a string which holds the password needed to connect to the SRRS
 */
const RefboxField = props => {
  return (
    <Refboxes>
      <Refbox>
        {tags(props).map(tag => {
          return (
            <RefboxButton
              isPenalty={props.isPenalty}
              key={tag}
              teamColor={"magenta"}
              id={`${tag}_magenta`}
              onClick={() => {
                console.log(`${TAG_VALUES[tag]}_magenta`);
                sendToServer(
                  `${TAG_VALUES[tag]}_magenta`,
                  "command",
                  props.session
                );
              }}
            >
              {tag}
            </RefboxButton>
          );
        })}
      </Refbox>
      <Refbox>
        {tags(props).map(tag => {
          return (
            <RefboxButton
              isPenalty={props.isPenalty}
              key={tag}
              teamColor={"cyan"}
              id={`${tag}_cyan`}
              onClick={() => {
                sendToServer(
                  `${TAG_VALUES[tag]}_cyan`,
                  "command",
                  props.session
                );
              }}
            >
              {tag}
            </RefboxButton>
          );
        })}
      </Refbox>
    </Refboxes>
  );
};

const Refboxes = styled.div`
  display: flex;
  justify-content: space-around;
  flex-wrap: no-wrap;
  padding-top: 2rem;

  ${media.sm`
    padding-top: 0;
  `};
`;

const Refbox = styled.div`
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  max-width: 12rem;

  ${media.xs`
    max-width: 16rem;
  `};
`;

const RefboxButton = styled(Button)`
  font-size: 2.5rem;

  ${media.xs`
    font-size: 4rem;
  `};

  ${props =>
    props.isPenalty
      ? css`
          flex-basis: 100%;
        `
      : css`
          flex-basis: 50%;
          max-width: 6rem;
          max-height: 6rem;

          ${media.xs`
            max-width: 8rem;
            max-height: 8rem;
            font-size: 4rem;
          `};
        `};
`;

const tags = props => {
  if (props.isPenalty) {
    return ["Soft", "Medium", "Hard"];
  } else {
    return ["KO", "FK", "GK", "TI", "C", "P"];
  }
};

RefboxField.propType = {
  isPenalty: PropTypes.bool.isRequired
};

export default RefboxField;
