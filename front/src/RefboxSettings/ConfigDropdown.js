import React from "react";
import styled from "styled-components";
import PropTypes from "prop-types";
import Dropdown from "../Dropdown";

const MODE_VALUES = ["Practice", "Dual team", "Match"];
const CONFIG_VALUES = [
  "Role assigner on",
  "Role assigner off",
  "Pass demo",
  "Penalty mode",
  "Ball Handling demo"
];

/**
 * Constructs two dropdowns for the type of game and role
 * Author: G.W. van der Heijden
 * Author: S.A. Tanja
 * Author: T.T.P. Franken
 *
 * Props:
 * - value: current value of the dropdown
 * - onChange: a function on what to do when the button is pressed
 */
const ConfigDropdown = props => {
  return (
    <DropDownSection>
      <RefboxDropdown
        currentValue={props.value}
        values={MODE_VALUES}
        onChange={value => {
          props.onChange("mode", value);
        }}
        enabled={true}
      />
      <RefboxDropdown
        currentValue={props.value}
        values={CONFIG_VALUES}
        onChange={value => {
          props.onChange("role", value);
        }}
        enabled={true}
      />
    </DropDownSection>
  );
};

const DropDownSection = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
  justify-content: center;
`;

const RefboxDropdown = styled(Dropdown)`
  width: 75%;
  margin-left: auto;
  margin-right: auto;
  height: 4rem;
`;

ConfigDropdown.propType = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired
};

export default ConfigDropdown;
