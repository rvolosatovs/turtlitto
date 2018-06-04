import React from "react";
import DropBall from "./DropBall";
import InOutButton from "./InOutButton";
import styled from "styled-components";
import PropTypes from "prop-types";

/**
 * Gives the settings part of the refbox:
 * - Drop ball
 * - Go in
 * - Go out
 *
 * Author: G.W. van der Heijden
 * Author: S.A. Tanja
 * Author: T.T.P. Franken
 *
 * Props:
 * - onClickDropBall: a function on what to do when DropBall is pressed
 * - onClickInOut: a function on what to do when InOutButton is pressed
 */
const RefboxSettings = props => {
  return (
    <Refbox>
      <DropBallButton
        onClick={() => {
          props.onClickDropBall();
        }}
      />
      <ButtonBlockWrapper>
        <InOutButton
          onClick={prop => {
            props.onClickInOut(prop);
          }}
        />
      </ButtonBlockWrapper>
    </Refbox>
  );
};

const ButtonBlockWrapper = styled.div`
  flex: 1;
  display: flex;
  flex-direction: row;
  justify-content: center;
`;

const DropBallButton = styled(DropBall)`
  margin: 0.5rem;
`;

const Refbox = styled.div`
  height: 50%;
  display: flex;
  flex-direction: column;
  justify-content: center;
`;

RefboxSettings.propType = {
  onClickDropBall: PropTypes.func.isRequired,
  onClickInOut: PropTypes.func.isRequired,
  onChange: PropTypes.func.isRequired
};

export default RefboxSettings;
