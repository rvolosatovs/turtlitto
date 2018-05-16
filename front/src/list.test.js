import React from "react";
import renderer from "react-test-renderer";
import list from "./list";

const TestComponent = props => {
  return <p>{props.data}</p>;
};

test("Multiple elements are rendered if enabled", () => {
  const component = renderer.create(
    list(TestComponent, [
      { enabled: true, data: 2, id: 2 },
      { enabled: true, data: 1, id: 1 }
    ])
  );
  expect(component.toJSON()).toMatchSnapshot();
});

test("No elements are added when no elements are present", () => {
  const component = renderer.create(list(TestComponent, []));
  expect(component.toJSON()).toMatchSnapshot();
});

test("Elements are only rendered when they are enabled", () => {
  const component = renderer.create(
    list(TestComponent, [
      { enabled: false, data: 2, id: 2 },
      { enabled: true, data: 1, id: 1 }
    ])
  );
  expect(component.toJSON()).toMatchSnapshot();
});
