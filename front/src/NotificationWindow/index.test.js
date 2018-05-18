import React from "react";
import { NotificationTypes, NotificationWindow } from ".";
import renderer from "react-test-renderer";

describe("NotificationWindow", () => {
  it("shows a SUCCESS notification", () => {
    const component = renderer.create(
      <NotificationWindow
        notificationType={NotificationTypes.SUCCESS}
        onDismiss={() => {
          console.log("Notification dismissed");
        }}
        message="This is a very annoying and extremely long error message which is displayed in full and hopefully uses wrap correctly with extremelylongwordswhichareuninterruptedandwithoutspaceswhichevenfillstheentirebox."
      />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
  it("shows a WARNING notification", () => {
    const component = renderer.create(
      <NotificationWindow
        notificationType={NotificationTypes.WARNING}
        onDismiss={() => {
          console.log("Notification dismissed");
        }}
        message="This is a very annoying and extremely long error message which is displayed in full and hopefully uses wrap correctly with extremelylongwordswhichareuninterruptedandwithoutspaceswhichevenfillstheentirebox."
      />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
  it("shows a ERROR notification", () => {
    const component = renderer.create(
      <NotificationWindow
        notificationType={NotificationTypes.ERROR}
        onDismiss={() => {
          console.log("Notification dismissed");
        }}
        message="This is a very annoying and extremely long error message which is displayed in full and hopefully uses wrap correctly with extremelylongwordswhichareuninterruptedandwithoutspaceswhichevenfillstheentirebox."
      />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
  it("shows a non identified notification", () => {
    const component = renderer.create(
      <NotificationWindow
        notificationType={NotificationTypes.SOMETHING}
        onDismiss={() => {
          console.log("Notification dismissed");
        }}
        message="This is a very annoying and extremely long error message which is displayed in full and hopefully uses wrap correctly with extremelylongwordswhichareuninterruptedandwithoutspaceswhichevenfillstheentirebox."
      />
    );
    let tree = component.toJSON();
    expect(tree).toMatchSnapshot();
  });
});
