import PropTypes from "prop-types";
import React from "react";
import Dropdown from "../Dropdown";
import Battery from "../Battery";
import styled from "styled-components";

const TEAM_VALUES = ["Magenta", "Cyan"];
const HOME_VALUES = ["Yellow home", "Blue home"];
const ROLE_VALUES = [
  "INACTIVE",
  "ROLE_NONE",
  "Att main",
  "Att assist",
  "Def main",
  "Def assist",
  "Def assist 2",
  "Goalkeeper"
];

/**
 * Show all details for a turtle
 * Author: H.E. van der Laan
 *
 * Props:
 * - turtle: An object containing the following turtle details:
 *   - battery: the current battery status of the turtle
 *   - home: the current home goal of this turtle
 *   - id: the ID number of the turtle
 *   - role: the current role of this turtle
 *   - team: the current team of this turtle
 * - editable: whether this turtle's properties can be edited
 * - onChange: a function with two arguments that is called when one of the dropdowns is changed. The first argument is name of the prop that is changed, the second argument is its new value.
 */
const Turtle = props => {
  const { battery, home, id, role, team } = props.turtle;
  const { editable, onChange } = props;
  return (
    <DefaultTurtle>
      <BatterySection>
        <Battery percentage={battery} />
      </BatterySection>
      <SubSection>
        <p>Turtle {id}</p>
      </SubSection>
      <DropDownSection>
        <Dropdown
          id={`turtle${id}__role`}
          currentValue={role}
          enabled={editable}
          values={ROLE_VALUES}
          onChange={value => {
            onChange("role", value);
          }}
        />
        <Dropdown
          id={`turtle${id}__home`}
          currentValue={home}
          enabled={editable}
          values={HOME_VALUES}
          onChange={value => {
            onChange("home", value);
          }}
        />
        <Dropdown
          id={`turtle${id}__team`}
          currentValue={team}
          enabled={editable}
          values={TEAM_VALUES}
          onChange={value => {
            onChange("team", value);
          }}
        />
      </DropDownSection>
    </DefaultTurtle>
  );
};

Turtle.propTypes = {
  turtle: PropTypes.shape({
    battery: PropTypes.number,
    home: PropTypes.string,
    id: PropTypes.number,
    role: PropTypes.string,
    team: PropTypes.string
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
