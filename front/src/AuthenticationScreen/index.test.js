import React from "react";
import { shallow } from "enzyme";
import AuthenticationScreen from ".";
import connectionTypes from "../BottomBar/connectionTypes";

describe("AuthenticationScreen", () => {
  describe("The user enters a incorrect token", () => {
    it("should match snapshot", () => {
      const onSubmit = (token, onIncorrectToken) => {
        onIncorrectToken();
      };
      const wrapper = shallow(
        <AuthenticationScreen
          onSubmit={onSubmit}
          connectionStatus={connectionTypes.CONNECTED}
        />
      );

      wrapper.find("#login-button").simulate("click");
      expect(wrapper).toMatchSnapshot();
    });
  });
});
