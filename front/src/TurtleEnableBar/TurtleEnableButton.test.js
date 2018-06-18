import React from "react";
import TurtleEnableButton from "./TurtleEnableButton";
import { shallowWithTheme } from "../testUtils";
import theme from "../theme";

//Test_items: TurtleEnableButton.js
//Input_spec: -
//Output_spec: -
//Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
describe("TurtleEnableButton", () => {
  describe("is active", () => {
    it("should match snapshot", () => {
      const wrapper = shallowWithTheme(
        <TurtleEnableButton enabled id={1} onTurtleEnableChange={() => {}} />
      ).dive();

      expect(wrapper).toMatchSnapshot();
      expect(wrapper).toHaveStyleRule("background-color", theme.buttonActive);
      expect(wrapper).toHaveStyleRule("border-color", "none");
      expect(wrapper).toHaveStyleRule(
        "border-style",
        theme.buttonBorderStyleActive
      );
    });
  });

  describe("is inactive", () => {
    it("should match snapshot", () => {
      const wrapper = shallowWithTheme(
        <TurtleEnableButton
          id={1}
          enabled={false}
          onTurtleEnableChange={() => {}}
        />
      ).dive();

      expect(wrapper).toMatchSnapshot();
      expect(wrapper).toHaveStyleRule("background-color", theme.button);
      expect(wrapper).toHaveStyleRule("border-color", theme.buttonBorder);
      expect(wrapper).toHaveStyleRule("border-style", theme.buttonBorderStyle);
    });
  });
});
