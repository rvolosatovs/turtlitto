import React from "react";
import TurtleEnableBar from ".";
import TurtleEnableButton from "./TurtleEnableButton";
import { shallowWithTheme, mountWithTheme } from "../testUtils";
import theme from "../theme";
import sinon from "sinon";

//Test_items: TurtleEnableBar index.js
//Input_spec: -
//Output_spec: -
//Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
describe("TurtleEnableBar", () => {
  it("should match snapshot", () => {
    const wrapper = shallowWithTheme(
      <TurtleEnableBar turtles={[]} onTurtleEnableChange={() => {}} />
    ).dive();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toHaveStyleRule("background", theme.turtleEnableBar);
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

  describe("the user enables a turtle", () => {
    it("should call `onClick` with correct turtle id", () => {
      const turtles = [
        {
          id: "1",
          enabled: false
        },
        {
          id: "2",
          enabled: false
        }
      ];
      const onTurtleEnableChangeSpy = sinon.spy();
      const wrapper = mountWithTheme(
        <TurtleEnableBar
          turtles={turtles}
          onTurtleEnableChange={onTurtleEnableChangeSpy}
        />
      );

      wrapper.find(TurtleEnableButton).forEach(button => {
        expect(button.props().enabled).toBe(false);
      });

      const buttonWrapper = wrapper.find(TurtleEnableButton).first();
      buttonWrapper.simulate("click");

      expect(onTurtleEnableChangeSpy.calledOnce).toBe(true);
      expect(onTurtleEnableChangeSpy.calledWithExactly("1")).toBe(true);
    });
  });
});
