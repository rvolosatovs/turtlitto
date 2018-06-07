import React from "react";
import { shallow } from "enzyme";
import TurtleList from ".";

describe("TurtleList", () => {
  it("should match snapshot", () => {
    const wrapper = shallow(
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
