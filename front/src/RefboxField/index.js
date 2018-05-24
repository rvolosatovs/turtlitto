import React from "react";
import RefboxButton from "./RefboxButton";
import styled from "styled-components";
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
 */
const RefboxField = props => {
  return (
    <div>
      <Refbox>
        {["KO", "FK", "GK", "TI", "C", "P"].map(tag => {
          return (
            <RefboxButton
              key={tag}
              teamColor={props.teamColor}
              onClick={() => {
                props.onClick(tag, props.teamColor);
              }}
              tag={tag}
            />
          );
        })}
      </Refbox>
    </div>
  );
};

/**
 * Places the buttons in 2 columns and rows of 3
 * Allows for multiple refboxes to be next to each other (float : left)
 */
const Refbox = styled.div`
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  margin: auto;
  width: 41.2rem;
  float: center;
`;

RefboxField.propType = {
  teamColor: PropTypes.oneOf(["cyan", "magenta"]).isRequired
};

export default RefboxField;
