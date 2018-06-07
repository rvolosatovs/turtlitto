import React from "react";
import ConnectionBar from ".";
import { shallow } from "enzyme";
import connectionTypes from "../connectionTypes";

describe("ConnectionBar", () => {
  describe("should match snapshot", () => {
    it("when in the connecting state", () => {
      const wrapper = shallow(
        <ConnectionBar connectionStatus={connectionTypes.CONNECTING} />
      );

      expect(wrapper).toMatchSnapshot();
    });
    it("when in the connected state", () => {
      const wrapper = shallow(
        <ConnectionBar connectionStatus={connectionTypes.CONNECTED} />
      );

      expect(wrapper).toMatchSnapshot();
    });
    it("when in the disconnected state", () => {
      const wrapper = shallow(
        <ConnectionBar connectionStatus={connectionTypes.DISCONNECTED} />
      );

      expect(wrapper).toMatchSnapshot();
    });
  });
});
