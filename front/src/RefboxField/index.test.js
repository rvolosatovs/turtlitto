import React from "react";
import RefboxField from ".";
import renderer from "react-test-renderer";

describe("RefboxField", () => {
  it("shows a cyan refbox field with all buttons", () => {
    const component = renderer.create(
      <RefboxField teamColor="cyan" onClick={(tag, color) => {}} />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
  it("shows a magenta refbox field with all buttons", () => {
    const component = renderer.create(
      <RefboxField teamColor="magenta" onClick={(tag, color) => {}} />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
