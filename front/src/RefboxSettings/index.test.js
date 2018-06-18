import React from "react";
import RefboxSettings from ".";
import { mountWithTheme } from "../testUtils";

describe("RefboxSettings", () => {
  it("shows the settings including two dropdown menus and three buttons", () => {
    const wrapper = mountWithTheme(<RefboxSettings />);
    expect(wrapper).toMatchSnapshot();
  });
});
