import React from "react";
import Button from "./RefboxButton";
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
 *  - onClick: a function on what to do when the button is pressed
 *  - isPenalty: a boolean which indicates whether to go into penalty mode
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
              onClick={() => {
                props.onClick(tag, "magenta");
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
              onClick={() => {
                props.onClick(tag, "cyan");
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
  flex-wrap: wrap;
  padding-top: 2rem;
`;

const Refbox = styled.div`
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  max-width: 12rem;

  /* TODO: define media query globally with sk */
  @media screen and (min-width: 360px) {
    max-width: 16rem;
  }
`;

const RefboxButton = styled(Button)`
  font-size: 2.5rem;

  /* TODO: define media query globally with sk */
  @media screen and (min-width: 360px) {
    font-size: 4rem;
  }

  ${props =>
    props.isPenalty
      ? css`
          flex-basis: 100%;
        `
      : css`
          flex-basis: 50%;
          max-width: 6rem;
          max-height: 6rem;

          /* TODO: define media query globally with sk */
          @media screen and (min-width: 360px) {
            max-width: 8rem;
            max-height: 8rem;
            font-size: 4rem;
          }
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
  isPenalty: PropTypes.bool.isRequired,
  onClick: PropTypes.func.isRequired
};

export default RefboxField;
