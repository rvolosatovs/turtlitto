import React from "react";
import NotificationWindow from ".";
import notificationTypes from "./notificationTypes";
import renderer from "react-test-renderer";

describe("NotificationWindow", () => {
  describe("is success notification", () => {
    it("should match snapshot", () => {
      const notification = {
        notificationType: notificationTypes.SUCCESS,
        message: "notification message"
      };
      const component = renderer.create(
        <NotificationWindow notification={notification} onDismiss={() => {}} />
      );
      let tree = component.toJSON();
      expect(tree).toMatchSnapshot();
    });
  });

  describe("is warning notification", () => {
    it("should match snapshot", () => {
      const notification = {
        notificationType: notificationTypes.WARNING,
        message: "notification message"
      };
      const component = renderer.create(
        <NotificationWindow notification={notification} onDismiss={() => {}} />
      );
      let tree = component.toJSON();
      expect(tree).toMatchSnapshot();
    });
  });

  describe("is error notification", () => {
    it("should match snapshot", () => {
      const notification = {
        notificationType: notificationTypes.ERROR,
        message: "notification message"
      };
      const component = renderer.create(
        <NotificationWindow notification={notification} onDismiss={() => {}} />
      );
      let tree = component.toJSON();
      expect(tree).toMatchSnapshot();
    });
  });

  describe("is null", () => {
    it("should not render anything", () => {
      const notification = null;
      const component = renderer.create(
        <NotificationWindow notification={notification} onDismiss={() => {}} />
      );
      let tree = component.toJSON();
      expect(tree).toMatchSnapshot();
    });
  });
});
