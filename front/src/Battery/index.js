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
  const { percentage, className } = props;
  return (
    <Section className={className}>
      {getBatteryIcon(percentage)}
      <CurrentBatteryIndicator>{percentage}%</CurrentBatteryIndicator>
    </Section>
  );
};

const Section = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const CurrentBatteryIndicator = styled.span`
  font-size: 4rem;
  padding-top: 0.6rem;
`;

const getBatteryIcon = percentage => {
  if (percentage < 10) {
    return (
      <FontAwesomeIcon
        icon={faBatteryEmpty}
        rotation={270}
        color="red"
        size="2x"
      />
    );
  } else if (percentage >= 10 && percentage < 40) {
    return <FontAwesomeIcon icon={faBatteryQuarter} rotation={270} size="2x" />;
  } else if (percentage >= 40 && percentage < 60) {
    return <FontAwesomeIcon icon={faBatteryHalf} rotation={270} size="2x" />;
  } else if (percentage >= 60 && percentage < 90) {
    return (
      <FontAwesomeIcon icon={faBatteryThreeQuarters} rotation={270} size="2x" />
    );
  } else {
    return <FontAwesomeIcon icon={faBatteryFull} rotation={270} size="2x" />;
  }
};

export default Battery;
