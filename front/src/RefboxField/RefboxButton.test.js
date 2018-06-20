import React from "react";
import RefboxButton from "./RefboxButton";
import { shallowWithTheme } from "../testUtils";
import theme from "../theme";

/*
 * Test_items: Refboxbutton.js
 * Input_spec: -
 * Output_spec: -
 * Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
 */
describe("RefboxButton", () => {
  describe("is cyan", () => {
    it("should match snapshot", () => {
      const wrapper = shallowWithTheme(
        <RefboxButton teamColor="cyan" onClick={() => {}} tag="KO" />
      ).dive();
      expect(wrapper).toMatchSnapshot();
      expect(wrapper).toHaveStyleRule("background-color", theme.refboxCyan);
      expect(wrapper).toHaveStyleRule("border-color", theme.refboxCyanBorder);
      expect(wrapper).toHaveStyleRule("border", "0.2rem solid");
      expect(wrapper).toHaveStyleRule(
        "background-color",
        theme.refboxCyanActive,
        { modifier: ":active" }
      );
    });
  });

  describe("is magenta", () => {
    it("should match snapshot", () => {
      const wrapper = shallowWithTheme(
        <RefboxButton teamColor="magenta" onClick={() => {}} tag="P" />
      ).dive();
      expect(wrapper).toMatchSnapshot();
      expect(wrapper).toHaveStyleRule("background-color", theme.refboxMagenta);
      expect(wrapper).toHaveStyleRule(
        "border-color",
        theme.refboxMagentaBorder
      );
      expect(wrapper).toHaveStyleRule("border", "0.2rem solid");
      expect(wrapper).toHaveStyleRule(
        "background-color",
        theme.refboxMagentaActive,
        { modifier: ":active" }
      );
    });
  });
});
