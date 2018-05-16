import React from "react";
import FontAwesomeIcon from "@fortawesome/react-fontawesome";
import faBatteryEmpty from "@fortawesome/fontawesome-free-solid/faBatteryEmpty";
import faBatteryQuarter from "@fortawesome/fontawesome-free-solid/faBatteryQuarter";
import faBatteryHalf from "@fortawesome/fontawesome-free-solid/faBatteryHalf";
import faBatteryThreeQuarters from "@fortawesome/fontawesome-free-solid/faBatteryThreeQuarters";
import faBatteryFull from "@fortawesome/fontawesome-free-solid/faBatteryFull";
import faAngleDown from "@fortawesome/fontawesome-free-solid/faAngleDown";

// Properties: percentage
const Battery = props => {
  const BatteryIcon = properties => {
    const percentage = properties.percentage;
    if (percentage < 10) {
      return (
        <FontAwesomeIcon icon={faBatteryEmpty} rotation={270} color="red" />
      );
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
  return (
    <div>
      {" "}
      <BatteryIcon icon={faBatteryFull} />
      {props.percentage}%
    </div>
  );
};

export default Battery;
