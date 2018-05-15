import React from "react";
import RSSButton from "./RSSButton";
import renderer from "react-test-renderer";

const testFunc = () => {
  return "Yeah, that'll do.";
};

describe("Button can be clicked", () => {
  it("Test Button with Snapshot", () => {
    const Button = renderer.create(
      <RSSButton buttonText={"hurrdurr"} onClick={testFunc()} enabled={true} />
    );
    let but = Button.toJSON();
    expect(but).toMatchSnapshot();
    const value = but.props.onClick;
    expect(value).toEqual("Yeah, that'll do.");
  });
});
