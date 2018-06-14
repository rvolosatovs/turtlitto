import React from "react";
import PropTypes from "prop-types";
import styled from "styled-components";
import Dropdown from "../Dropdown";
import TurtleList from "../TurtleList";
import sendToServer from "../sendToServer";

/**
 * Show the settings of all turtles.
 * Author: H.E. van der Laan
 * Author: T.T.P. Franken
 *
 * props:
 * - turtles: an array of Turtles
 * - session: a string which holds the password needed to connect to the SRRS
 */

const CONFIG_VALUES = [
  "Role assigner on",
  "Role assigner off",
  "Pass demo",
  "Penalty mode",
  "Ball Handling demo"
];

const COMMAND_VALUES = {
  "Role assigner on": "role_assigner_on",
  "Role assigner off": "role_assigner_off",
  "Pass demo": "pass_demo",
  "Penalty mode": "penalty_demo",
  "Ball Handling demo": "ball_handling_demo"
};

const Settings = props => {
  const { turtles, session } = props;
  return (
    <SettingsWrapper>
      <RoleDropdown
        id={"settings_role-dropdown"}
        currentValue={"Whatever"}
        values={CONFIG_VALUES}
        onChange={value => {
          sendToServer(COMMAND_VALUES[value], "command", session);
        }}
        enabled={true}
      />
      <TurtleList turtles={turtles} session={session} />
    </SettingsWrapper>
  );
};

Settings.propTypes = {
  turtles: PropTypes.objectOf(
    PropTypes.shape({
      enabled: PropTypes.bool.isRequired,
      batteryvoltage: PropTypes.number.isRequired,
      homegoal: PropTypes.string.isRequired,
      role: PropTypes.string.isRequired,
      teamcolor: PropTypes.string.isRequired
    })
  ).isRequired
};

const RoleDropdown = styled(Dropdown)`
  margin-left: auto;
  margin-right: auto;
  height: 4rem;
`;

const SettingsWrapper = styled.div`
  display: flex;
  flex-direction: column;
`;

export default Settings;
