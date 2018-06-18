import React from "react";
import ConnectionBar from ".";
import connectionTypes from "../connectionTypes";
import { mountWithTheme, shallowWithTheme } from "../../testUtils";
import theme from "../../theme";

//Test_items: ConnectionBar index.js
//Input_spec: -
//Output_spec: -
//Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
describe("ConnectionBar", () => {
  describe("should match snapshot", () => {
    it("when in the connecting state", () => {
      const wrapper = mountWithTheme(
        <ConnectionBar connectionStatus={connectionTypes.CONNECTING} />
      );

      expect(wrapper).toMatchSnapshot();
      expect(wrapper).toHaveStyleRule("background", theme.warning);
    });
    it("when in the connected state", () => {
      const wrapper = mountWithTheme(
        <ConnectionBar connectionStatus={connectionTypes.CONNECTED} />
      );

      expect(wrapper).toMatchSnapshot();
      expect(wrapper).toHaveStyleRule("background", theme.success);
    });
    it("when in the disconnected state", () => {
      const wrapper = mountWithTheme(
        <ConnectionBar connectionStatus={connectionTypes.DISCONNECTED} />
      );

      expect(wrapper).toMatchSnapshot();
      expect(wrapper).toHaveStyleRule("background", theme.error);
    });
  });

  it("should throw an error when an unknown connection type is passed", () => {
    const connectionType = "unknown-connection-type";
    const oldError = console.error;

    console.error = () => {};
    expect(() => {
      shallowWithTheme(<ConnectionBar connectionStatus={connectionType} />);
    }).toThrow();
    console.error = oldError;
  });
});
