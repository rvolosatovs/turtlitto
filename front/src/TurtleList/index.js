import React from "react";
import PropTypes from "prop-types";
import { Grid, Row, Col } from "react-flexbox-grid";
import Turtle from "../Turtle";

/**
 * List containing for each turtle their respective status.
 * Author: B. Afonins
 * Author: H.E. van der Laan
 *
 * Props:
 *  - turtles: an object containing the following turtle details:
 *   - enabled: whether the turtle is enabled
 *   - batteryvoltage: the current battery status of the turtle
 *   - homegoal: the current home goal of this turtle
 *   - role: the current role of this turtle
 *   - teamcolor: the current team of this turtle
 */
const TurtleList = props => {
  const { turtles, session } = props;
  return (
    <Grid>
      <Row>
        {Object.keys(turtles)
          .filter(id => turtles[id].enabled)
          .map((id, index) => {
            return (
              <Col xs={12} md={6} key={index}>
                <Turtle
                  key={index}
                  id={id}
                  turtle={turtles[id]}
                  session={session}
                  editable
                />
              </Col>
            );
          })}
      </Row>
    </Grid>
  );
};

TurtleList.propTypes = {
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

export default TurtleList;
