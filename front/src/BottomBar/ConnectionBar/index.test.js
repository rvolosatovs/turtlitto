import React from "react";
import ConnectionBar from ".";
import connectionTypes from "../connectionTypes";
import { mountWithTheme } from "../../testUtils";

describe("ConnectionBar", () => {
  describe("should match snapshot", () => {
    it("when in the connecting state", () => {
      const wrapper = mountWithTheme(
        <ConnectionBar connectionStatus={connectionTypes.CONNECTING} />
      );

      expect(wrapper).toMatchSnapshot();
    });
    it("when in the connected state", () => {
      const wrapper = mountWithTheme(
        <ConnectionBar connectionStatus={connectionTypes.CONNECTED} />
      );

      expect(wrapper).toMatchSnapshot();
    });
    it("when in the disconnected state", () => {
      const wrapper = mountWithTheme(
        <ConnectionBar connectionStatus={connectionTypes.DISCONNECTED} />
      );

      expect(wrapper).toMatchSnapshot();
    });
  });
});
