import React from "react";
import RefboxSettings from ".";
import { mountWithTheme } from "../testUtils";

/*
 * Test_items: Settings index.js
 * Input_spec: -
 * Output_spec: -
 * Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
 */
describe("RefboxSettings", () => {
  it("shows the settings including two dropdown menus and three buttons", () => {
    const wrapper = mountWithTheme(<RefboxSettings />);
    expect(wrapper).toMatchSnapshot();
  });
});
