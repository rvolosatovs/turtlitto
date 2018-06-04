import React from "react";
import NotificationWindow from ".";
import notificationTypes from "./notificationTypes";
import { shallow } from "enzyme";

describe("NotificationWindow", () => {
  it("should display a success notification", () => {
    const notification = {
      notificationType: notificationTypes.SUCCESS,
      message: "notification message"
    };
    const wrapper = shallow(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
  });

  it("should display a warning notification", () => {
    const notification = {
      notificationType: notificationTypes.WARNING,
      message: "notification message"
    };
    const wrapper = shallow(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
  });

  it("should display an error notification", () => {
    const notification = {
      notificationType: notificationTypes.ERROR,
      message: "notification message"
    };
    const wrapper = shallow(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
  });

  it("should not render anything", () => {
    const notification = null;
    const wrapper = shallow(
      <NotificationWindow notification={notification} onDismiss={() => {}} />
    );
    expect(wrapper).toMatchSnapshot();
  });

  describe("gets invalid notificationType", () => {
    it("should throw an error", () => {
      const notification = "invalid";
      expect(() => {
        shallow(
          <NotificationWindow
            notification={notification}
            onDismiss={() => {}}
          />
        );
      }).toThrow();
    });
  });
});
