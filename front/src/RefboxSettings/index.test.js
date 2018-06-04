import React from "react";
import RefboxSettings from ".";
import { shallow } from "enzyme";

describe("RefboxSettings", () => {
  it("shows the settings including two dropdown menus and three buttons", () => {
    const wrapper = shallow(<RefboxSettings />);
    expect(wrapper).toMatchSnapshot();
  });
});
