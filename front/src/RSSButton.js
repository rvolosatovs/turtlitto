import PropTypes from "prop-types";
import React from "react";

/**
 * Create a standard RSS Button
 *
 * Props:
 *
 */

const RSSButton = props => {
  const { buttonText, onClick, enabled } = props;
  return (
    <button onClick={onClick} disabled={!enabled}>
      {buttonText}
    </button>
  );
};

RSSButton.propTypes = {
  buttonText: PropTypes.string.isRequired,
  onClick: PropTypes.func.isRequired,
  enabled: PropTypes.bool.isRequired
};

export default RSSButton;
