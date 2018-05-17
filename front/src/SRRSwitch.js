import PropTypes from "prop-types";
import React from "react";

/**
 * Create a standard RSS Button
 *
 * Props:
 *
 */

const SRRSwitch = props => {
  const { currentValue, onChange, enabled } = props;
  return (
    <input
      type="checkbox"
      onChange={onChange}
      disabled={!enabled}
      defaultChecked={currentValue}
    />
  );
};

SRRSwitch.propTypes = {
  currentValue: PropTypes.bool.isRequired,
  buttonText: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired,
  enabled: PropTypes.bool.isRequired
};

export default SRRSwitch;
