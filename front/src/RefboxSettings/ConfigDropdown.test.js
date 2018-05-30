import React from "react";
import ConfigDropdown from "./ConfigDropdown";
import { shallow } from "enzyme";

describe("ConfigDropdown", () => {
  it("shows two dropdown menus", () => {
    const wrapper = shallow(<ConfigDropdown value="" onChange={() => {}} />);
    expect(wrapper).toMatchSnapshot();
  });
});
