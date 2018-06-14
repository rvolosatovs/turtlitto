import React from "react";
import DropBall from "./DropBall";
import InOutButton from "./InOutButton";
import styled from "styled-components";
import sendToServer from "../sendToServer";

/**
 * Gives the settings part of the refbox:
 * - Drop ball
 * - Go in
 * - Go out
 *
 * Author: G.W. van der Heijden
 * Author: S.A. Tanja
 * Author: T.T.P. Franken
 */
const RefboxSettings = props => {
  return (
    <Refbox className={props.className}>
      <DropBallButton
        onClick={() => {
          sendToServer("dropped_ball", "command", props.session);
        }}
      />
      <ButtonBlockWrapper>
        <InOutButton
          onClick={prop => {
            sendToServer(prop, "command", props.session);
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
  padding: 0.5rem;
`;

export default RefboxSettings;
