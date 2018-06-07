import React from "react";
import PropTypes from "prop-types";
import { Row, Col } from "react-flexbox-grid";

import Turtle from "../Turtle";

const TurtleList = props => {
  const { turtles } = props;
  return (
    <Row>
      {Object.keys(turtles)
        .filter(id => turtles[id].enabled)
        .map((id, index) => {
          return (
            <Col xs={12} md={6} key={index}>
              <Turtle key={index} id={id} turtle={turtles[id]} editable />
            </Col>
          );
        })}
    </Row>
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
