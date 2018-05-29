import React from "react";
import RefboxSettings from ".";
import renderer from "react-test-renderer";

describe("RefboxSettings", () => {
  it("shows the settings including two dropdown menus and three buttons", () => {
    const component = renderer.create(<RefboxSettings />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
