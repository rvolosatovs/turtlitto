import React from "react";
import Dropdown from ".";
import { shallow } from "enzyme";

/* 
 * Test_items: Dropdown index.js
 * Input_spec: -
 * Output_spec: -
 * Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
 */
describe("Dropdown", () => {
  it("can be changed", () => {
    let value = "test";
    const wrapper = shallow(
      <Dropdown
        currentValue={value}
        values={["test", "test2"]}
        enabled
        onChange={newValue => {
          value = newValue;
        }}
      />
    );
    expect(wrapper).toMatchSnapshot();

    // Check if we can trigger onChange
    expect(value).toBe("test");
    wrapper.find("select").simulate("change", { target: { value: "test2" } });
    expect(value).toBe("test2");
  });
});
