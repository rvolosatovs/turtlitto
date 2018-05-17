import React from "react";
import Battery from "./Battery";
import renderer from "react-test-renderer";

test("Battery 100%", () => {
  const component = renderer.create(<Battery percentage={100} />);
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});

test("Battery 75%", () => {
  const component = renderer.create(<Battery percentage={75} />);
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});

test("Battery 50%", () => {
  const component = renderer.create(<Battery percentage={50} />);
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});

test("Battery 25%", () => {
  const component = renderer.create(<Battery percentage={25} />);
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});

test("Battery 0%", () => {
  const component = renderer.create(<Battery percentage={0} />);
  let tree = component.toJSON();
  expect(tree).toMatchSnapshot();
});
