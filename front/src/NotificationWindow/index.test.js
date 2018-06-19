import React from "react";
import NotificationWindow from ".";
import notificationTypes from "./notificationTypes";
import { shallow } from "enzyme";
import { mountWithTheme } from "../testUtils";
import theme from "../theme";

/* 
 * Test_items: Notification Window index.js
 * Input_spec: -
 * Output_spec: -
 * Envir_needs: snapshot (automatically made, found in the __snapshot__ folder).
 */
describe("NotificationWindow", () => {
  it("should display a success notification", () => {
    const notification = {
      notificationType: notificationTypes.SUCCESS,
      message: "notification message"
    };
    const wrapper = mountWithTheme(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toHaveStyleRule("background", theme.success);
  });

  it("should display a warning notification", () => {
    const notification = {
      notificationType: notificationTypes.WARNING,
      message: "notification message"
    };
    const wrapper = mountWithTheme(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toHaveStyleRule("background", theme.warning);
  });

  it("should display an error notification", () => {
    const notification = {
      notificationType: notificationTypes.ERROR,
      message: "notification message"
    };
    const wrapper = mountWithTheme(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toHaveStyleRule("background", theme.error);
  });

  it("should not render anything", () => {
    const notification = null;
    const wrapper = shallow(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
  });

  it("throws an error when an unknown notification type is shown", () => {
    const notification = {
      notificationType: "annoying",
      message: "Crash pls"
    };
    // Hide PropTypes validation error, we want the normal error
    const oldError = console.error;
    console.error = () => {};
    expect(() => {
      shallow(
        <NotificationWindow notification={notification} onDismiss={() => {}} />
      );
    }).toThrow();
    console.error = oldError;
  });

  it("can be dismissed", () => {
    const dismissFn = jest.fn();
    const notification = {
      notificationType: notificationTypes.ERROR,
      message: "Pants on FIRE"
    };
    const wrapper = shallow(
      <NotificationWindow
        notification={notification}
        onDismiss={() => {
          dismissFn();
        }}
      />
    );
    expect(dismissFn).not.toBeCalled();
    wrapper
      .find("FontAwesomeIcon")
      .parent() // Needed because under shallow, events don't bubble
      .simulate("click");
    expect(dismissFn).toBeCalled();
  });
});
