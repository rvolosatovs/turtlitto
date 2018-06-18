import React from "react";
import TurtleList from ".";
import { shallowWithTheme, mountWithTheme } from "../testUtils";

//Test_items: TurtleList index.js
//Input_spec: -
//Output_spec: -
//Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
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
            homegoal: "yellow",
            role: "inactive",
            teamcolor: "magenta"
          },
          2: {
            enabled: false,
            batteryvoltage: 42,
            homegoal: "yellow",
            role: "inactive",
            teamcolor: "magenta"
          }
        }}
      />
    );
    expect(wrapper).toMatchSnapshot();
  });
});
