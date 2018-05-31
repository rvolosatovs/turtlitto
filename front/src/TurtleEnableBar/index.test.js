import React from "react";
import { shallow, mount } from "enzyme";
import TurtleEnableBar from ".";
import TurtleEnableButton from "./TurtleEnableButton";

describe("TurtleEnableBar", () => {
  it("should match snapshot", () => {
    const wrapper = shallow(
      <TurtleEnableBar turtles={[]} onTurtleEnableChange={() => {}} />
    );

    expect(wrapper).toMatchSnapshot();
  });

  it("should render all turtles as <TurtleEnableButton /> components", () => {
    const turtles = [
      {
        id: 0,
        enabled: false
      },
      {
        id: 1,
        enabled: false
      },
      {
        id: 2,
        enabled: false
      }
    ];
    const wrapper = mount(
      <TurtleEnableBar turtles={turtles} onTurtleEnableChange={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();

    expect(wrapper.find(TurtleEnableButton).length).toBe(turtles.length);
  });
});
