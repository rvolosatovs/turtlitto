import React from "react";
import Battery from ".";
import renderer from "react-test-renderer";

describe("Battery", () => {
  it("shows a full battery on 100%", () => {
    const component = renderer.create(<Battery percentage={100} />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });

  it("shows a 3/4 full battery on 75%", () => {
    const component = renderer.create(<Battery percentage={75} />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });

  it("shows a half full battery on 50%", () => {
    const component = renderer.create(<Battery percentage={50} />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });

  it("shows a 1/4 full battery on 25%", () => {
    const component = renderer.create(<Battery percentage={25} />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });

  it("shows an empty battery on 0%", () => {
    const component = renderer.create(<Battery percentage={0} />);
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
