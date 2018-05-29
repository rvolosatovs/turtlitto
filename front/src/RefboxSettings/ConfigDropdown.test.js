import React from "react";
import ConfigDropdown from "./ConfigDropdown";
import renderer from "react-test-renderer";

describe("ConfigDropdown", () => {
  it("shows two dropdown menus", () => {
    const component = renderer.create(
      <ConfigDropdown value="" onChange={() => {}} />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
