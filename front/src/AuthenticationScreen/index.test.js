import React from "react";
import AuthenticationScreen from ".";
import connectionTypes from "../BottomBar/connectionTypes";
import { mountWithTheme } from "../testUtils";

describe("AuthenticationScreen", () => {
  describe("The user enters a incorrect token", () => {
    it("should match snapshot", () => {
      const onSubmit = token => {};
      const wrapper = mountWithTheme(
        <AuthenticationScreen
          onSubmit={onSubmit}
          connectionStatus={connectionTypes.CONNECTED}
          notification={"hurdur"}
        />
      );

      wrapper
        .find("#login-button")
        .hostNodes()
        .simulate("click");
      expect(wrapper).toMatchSnapshot();
    });
  });
});
