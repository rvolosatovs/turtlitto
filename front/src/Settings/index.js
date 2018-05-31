import React from "react";
import PropTypes from "prop-types";
import Turtle from "../Turtle";

/**
 * Show the settings of all turtles.
 * Author: H.E. van der Laan
 *
 * props:
 * - turtles: an array of Turtles
 * - onTurtleEnableChange: function to call when the turtle enable button is pressed
 */
const Settings = props => {
  const { turtles } = props;
  return (
    <div>
      {turtles.map(turtle => (
        <Turtle
          key={turtle.id}
          turtle={turtle}
          editable
          onChange={(changedProp, newValue) => {} /* TODO: turtle update */}
        />
      ))}
    </div>
  );
};

Settings.propTypes = {
  turtles: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number.isRequired,
      enabled: PropTypes.bool.isRequired,
      battery: PropTypes.number.isRequired,
      home: PropTypes.string.isRequired,
      role: PropTypes.string.isRequired,
      team: PropTypes.string.isRequired
    })
  ).isRequired
};

export default Settings;
