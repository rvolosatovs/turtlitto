import React from "react";
import DropBall from "./DropBall";
import renderer from "react-test-renderer";

describe("DropBall", () => {
  it("shows a button with DB written on it", () => {
    const component = renderer.create(<DropBall onClick={() => {}} />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
