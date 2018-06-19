import React from "react";
import AuthenticationScreen from ".";
import connectionTypes from "../BottomBar/connectionTypes";
import { mountWithTheme } from "../testUtils";
import theme from "../theme";

describe("AuthenticationScreen", () => {
  it("should match snapshot", () => {
    const wrapper = mountWithTheme(
      <AuthenticationScreen
        onSubmit={() => {}}
        connectionStatus={connectionTypes.DISCONNECTED}
        notification={""}
      />
    );

    expect(wrapper).toMatchSnapshot();
  });

  it("should show warning label when notification passed", () => {
    const notificationText = "Incorrect token";
    const wrapper = mountWithTheme(
      <AuthenticationScreen
        onSubmit={() => {}}
        connectionStatus={connectionTypes.DISCONNECTED}
        notification={notificationText}
      />
    );

    const label = wrapper.find("WarningLabel");
    expect(label.text()).toEqual(notificationText);
    expect(label).toHaveStyleRule("color", theme.error);
  });
});
