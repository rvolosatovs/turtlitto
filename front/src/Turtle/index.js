import PropTypes from "prop-types";
import React from "react";
import Dropdown from "../Dropdown";
import Battery from "../Battery";
import styled from "styled-components";

// The following values are used in the Turtle refbox. The property name (or key) is shown in the UI, while its value will be sent over to the TRC.
const TEAM_VALUES = { Magenta: "magenta", Cyan: "cyan" };
const HOME_VALUES = { "Yellow home": "yellow", "Blue home": "blue" };
const ROLE_VALUES = {
  INACTIVE: "inactive",
  ROLE_NONE: "none",
  "Att main": "attacker_main",
  "Att assist": "attacker_assist",
  "Def main": "defender_main",
  "Def assist": "defender_assist",
  "Def assist 2": "defender_assist2",
  Goalkeeper: "goalkeeper"
};

const onChange = (id, propName, propValue) => {
  const body = {};
  body[id] = {};
  body[id][propName] = propValue;
  fetch("/api/v1/turtles/", {
    body: JSON.stringify(body),
    method: "PUT"
  })
    .then(data => console.log(data))
    .catch(error => console.error(error));
};

/**
 * Show all details for a turtle
 * Author: H.E. van der Laan
 *
 * Props:
 * - turtle: An object containing the following turtle details:
 *   - batteryvoltage: the current battery status of the turtle
 *   - homegoal: the current home goal of this turtle
 *   - id: the ID number of the turtle
 *   - role: the current role of this turtle
 *   - teamcolor: the current team of this turtle
 * - editable: whether this turtle's properties can be edited
 * - onChange: a function with two arguments that is called when one of the dropdowns is changed. The first argument is name of the prop that is changed, the second argument is its new value.
 */
const Turtle = props => {
  const { batteryvoltage, homegoal, role, teamcolor, id } = props.turtle;
  const { editable } = props;
  return (
    <DefaultTurtle>
      <BatterySection>
        <Battery percentage={batteryvoltage} />
      </BatterySection>
      <SubSection>
        <p>Turtle {id}</p>
      </SubSection>
      <DropDownSection>
        <Dropdown
          id={`turtle${id}__role`}
          currentValue={role}
          enabled={editable}
          values={Object.keys(ROLE_VALUES)}
          onChange={value => {
            onChange(id, "role", ROLE_VALUES[value]);
          }}
        />
        <Dropdown
          id={`turtle${id}__home`}
          currentValue={homegoal}
          enabled={editable}
          values={Object.keys(HOME_VALUES)}
          onChange={value => {
            onChange(id, "homegoal", HOME_VALUES[value]);
          }}
        />
        <Dropdown
          id={`turtle${id}__team`}
          currentValue={teamcolor}
          enabled={editable}
          values={Object.keys(TEAM_VALUES)}
          onChange={value => {
            onChange(id, "teamcolor", TEAM_VALUES[value]);
          }}
        />
      </DropDownSection>
    </DefaultTurtle>
  );
};

Turtle.propTypes = {
  turtle: PropTypes.shape({
    batteryvoltage: PropTypes.number,
    homegoal: PropTypes.string,
    role: PropTypes.string,
    teamcolor: PropTypes.string,
    id: PropTypes.string.isRequired
  }).isRequired,
  editable: PropTypes.bool,
  onChange: PropTypes.func
};

const BatterySection = styled.div`
  flex-basis: 35%;
`;

const SubSection = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  flex-basis: 35%;
`;

const DropDownSection = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex-basis: 30%;
`;

const DefaultTurtle = styled.div`
  width: 90%;
  margin: 0 auto;
  padding: 1.5rem 0;
  border-style: solid;
  border: 0.25rem;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  font-size: 2rem;
`;

export default Turtle;
