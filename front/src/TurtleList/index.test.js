import React from "react";
import TurtleList from ".";
import { shallowWithTheme, mountWithTheme } from "../testUtils";

describe("TurtleList", () => {
  it("should match snapshot", () => {
    const wrapper = shallowWithTheme(<TurtleList turtles={{}} />).dive();
    expect(wrapper).toMatchSnapshot();
  });

  it("should render all turtles", () => {
    const wrapper = mountWithTheme(
      <TurtleList
        turtles={{
          1: {
            enabled: true,
            batteryvoltage: 66,
            homegoal: "Yellow home",
            role: "INACTIVE",
            teamcolor: "Magenta"
          },
          2: {
            enabled: false,
            batteryvoltage: 42,
            homegoal: "Yellow home",
            role: "INACTIVE",
            teamcolor: "Magenta"
          }
        }}
      />
    );
    expect(wrapper).toMatchSnapshot();
  });
});
