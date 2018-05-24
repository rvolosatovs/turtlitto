import React from "react";
import NotificationWindow from ".";
import notificationTypes from "./notificationTypes";
import renderer from "react-test-renderer";

describe("NotificationWindow", () => {
  describe("is success notification", () => {
    it("should match snapshot", () => {
      const component = renderer.create(
        <NotificationWindow
          notificationType={notificationTypes.SUCCESS}
          onDismiss={() => {}}
          message="notification message"
        />
      );
      let tree = component.toJSON();
      expect(tree).toMatchSnapshot();
    });
  });

  describe("is warning notification", () => {
    it("should match snapshot", () => {
      const component = renderer.create(
        <NotificationWindow
          notificationType={notificationTypes.WARNING}
          onDismiss={() => {}}
          message="notification message"
        />
      );
      let tree = component.toJSON();
      expect(tree).toMatchSnapshot();
    });
  });

  describe("is error notification", () => {
    it("should match snapshot", () => {
      const component = renderer.create(
        <NotificationWindow
          notificationType={notificationTypes.ERROR}
          onDismiss={() => {}}
          message="notification message"
        />
      );
      let tree = component.toJSON();
      expect(tree).toMatchSnapshot();
    });
  });
});
