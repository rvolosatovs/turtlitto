import React from "react";
import InOutButton from "./InOutButton";
import renderer from "react-test-renderer";

describe("InOutButton", () => {
  it("shows two buttons with go in or go out written on it", () => {
    const component = renderer.create(<InOutButton onClick={() => {}} />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
