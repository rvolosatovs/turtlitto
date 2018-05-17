import PropTypes from "prop-types";
import React from "react";
import Dropdown from "./Dropdown";
import Battery from "./Battery";

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
 *   - editable: whether this turtle's properties can be edited
 *   - home: the current home goal of this robot
 *   - id: the ID number of the turtle
 *   - role: the current role of this robot
 *   - team: the current team of this robot
 */
const Turtle = props => {
  const { battery, editable, home, id, role, team } = props.turtle;
  return (
    <div className={props.className}>
      <div>
        <Battery percentage={battery} />
      </div>
      <div>
        <p>Turtle {id}</p>
      </div>
      <div>
        <Dropdown currentValue={role} enabled={editable} values={ROLE_VALUES} />
        <Dropdown currentValue={home} enabled={editable} values={HOME_VALUES} />
        <Dropdown currentValue={team} enabled={editable} values={TEAM_VALUES} />
      </div>
    </div>
  );
};

Turtle.propTypes = {
  turtle: PropTypes.shape({
    battery: PropTypes.number,
    editable: PropTypes.bool,
    home: PropTypes.string,
    id: PropTypes.number,
    role: PropTypes.string,
    team: PropTypes.string
  }).isRequired
};

export default Turtle;
