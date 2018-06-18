import React from "react";
import TurtleEnableBar from ".";
import TurtleEnableButton from "./TurtleEnableButton";
import { shallowWithTheme, mountWithTheme } from "../testUtils";

describe("TurtleEnableBar", () => {
  it("should match snapshot", () => {
    const wrapper = shallowWithTheme(
      <TurtleEnableBar turtles={[]} onTurtleEnableChange={() => {}} />
    ).dive();

    expect(wrapper).toMatchSnapshot();
  });

  it("should render all turtles as <TurtleEnableButton /> components", () => {
    const turtles = [
      {
        id: "1",
        enabled: false
      },
      {
        id: "2",
        enabled: false
      },
      {
        id: "3",
        enabled: false
      }
    ];
    const wrapper = mountWithTheme(
      <TurtleEnableBar turtles={turtles} onTurtleEnableChange={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();

    expect(wrapper.find(TurtleEnableButton).length).toBe(turtles.length);
  });
});
