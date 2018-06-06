import React from "react";
import PropTypes from "prop-types";
import { Row, Col } from "react-flexbox-grid";

import Turtle from "../Turtle";

const TurtleList = props => {
  const { turtles } = props;
  return (
    <Row>
      {turtles.map((turtle, index) => (
        <Col xs={12} md={6} key={index}>
          <Turtle key={index} turtle={turtle} editable />
        </Col>
      ))}
    </Row>
  );
};

TurtleList.propTypes = {
  turtles: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      enabled: PropTypes.bool.isRequired,
      batteryvoltage: PropTypes.number.isRequired,
      homegoal: PropTypes.string.isRequired,
      role: PropTypes.string.isRequired,
      teamcolor: PropTypes.string.isRequired
    })
  ).isRequired
};

export default TurtleList;
