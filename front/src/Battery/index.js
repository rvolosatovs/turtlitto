import React from "react";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import faBatteryEmpty from "@fortawesome/fontawesome-free-solid/faBatteryEmpty";
import faBatteryQuarter from "@fortawesome/fontawesome-free-solid/faBatteryQuarter";
import faBatteryHalf from "@fortawesome/fontawesome-free-solid/faBatteryHalf";
import faBatteryThreeQuarters from "@fortawesome/fontawesome-free-solid/faBatteryThreeQuarters";
import faBatteryFull from "@fortawesome/fontawesome-free-solid/faBatteryFull";
import styled from "styled-components";

/**
 * Show a nice battery icon
 * Author: G.W. van der Heijden
 * Author: H.E. van der Laan
 *
 * Props:
 * - percentage: battery percentage
 */
const Battery = props => {
  const percentage = props.percentage;
  return (
    <Section>
      {getBatteryIcon(percentage)}
      {percentage}%
    </Section>
  );
};

const Section = styled.div`
  flex: 1;
  font-size: 4rem;
  margin: auto;
`;

const getBatteryIcon = percentage => {
  if (percentage < 10) {
    return <FontAwesomeIcon icon={faBatteryEmpty} rotation={270} color="red" />;
  } else if (percentage >= 10 && percentage < 40) {
    return <FontAwesomeIcon icon={faBatteryQuarter} rotation={270} />;
  } else if (percentage >= 40 && percentage < 60) {
    return <FontAwesomeIcon icon={faBatteryHalf} rotation={270} />;
  } else if (percentage >= 60 && percentage < 90) {
    return <FontAwesomeIcon icon={faBatteryThreeQuarters} rotation={270} />;
  } else {
    return <FontAwesomeIcon icon={faBatteryFull} rotation={270} />;
  }
};

export default Battery;
