import PropTypes from "prop-types";
import React from "react";

/**
 * Create a dropdown menu
 *
 * Props:
 *  - currentValue: the current value of this dropdown
 *  - values: the possible values of this dropdown
 *  - onChange: a function to call when an update is posted. If not provided, the dropdown menu becomes readonly
 *  - enabled: if the dropdown should be enabled
 */
const Dropdown = props => {
  const { currentValue, values, onChange, enabled } = props;
  return (
    <select value={currentValue} onChange={onChange} disabled={!enabled}>
      {values.map(value => {
        return <option key={value}>{value}</option>;
      })}
    </select>
  );
};

Dropdown.propTypes = {
  currentValue: PropTypes.string.isRequired,
  values: PropTypes.arrayOf(PropTypes.string).isRequired,
  onChange: PropTypes.func.isRequired,
  enabled: PropTypes.bool.isRequired
};

export default Dropdown;
