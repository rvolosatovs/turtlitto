import PropTypes from "prop-types";
import React from "react";

/**
 * Create a standard RSS Button
 *
 * Props:
 *
 */

const SRRButton = props => {
  const { buttonText, onClick, enabled } = props;
  return (
    <button className={props.className} onClick={onClick} disabled={!enabled}>
      {buttonText}
    </button>
  );
};

SRRButton.propTypes = {
  buttonText: PropTypes.element.isRequired,
  onClick: PropTypes.func.isRequired,
  enabled: PropTypes.bool.isRequired
};

export default SRRButton;
