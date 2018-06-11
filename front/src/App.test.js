import React from "react";
import ReactDOM from "react-dom";
import { Server } from "mock-socket";
import { shallow } from "enzyme";
import App from "./App";

jest.useFakeTimers();

describe("App.js", () => {
  it("automatically reconnects", () => {
    //Establish a server and wait for a connection
    const mockServer = new Server("ws://user:testtoken@localhost/api/v1/state");
    mockServer.on("connection", server => {
      connectionCount++;
    });
    mockServer.on("close", server => {
      connectionCount--;
    });
    const wrapper = shallow(<App />);
    wrapper.setState({ token: "testtoken" });
    let connectionCount = 0;
    //Let the reconnect time run down
    jest.runAllTimers();
    //Check if exactly one connection has been made
    expect(connectionCount).toBe(1);
    //Close the server to disconnect and check if there is no connection
    mockServer.close();
    expect(connectionCount).toBe(0);
    //Establish a new server and wait for a connection
    const mockServer2 = new Server(
      "ws://user:testtoken@localhost/api/v1/state"
    );
    mockServer2.on("connection", server => {
      connectionCount++;
    });
    mockServer2.on("close", server => {
      connectionCount--;
    });
    jest.runAllTimers();
    //Wait one second since it will take one second to reconnect
    setTimeout(() => {
      //Check if exactly one connection has been made
      expect(connectionCount).toBe(1);
      wrapper.unmount(this);
    }, 1000);
    mockServer2.close();
  });

  it("renders without crashing", () => {
    const div = document.createElement("div");
    ReactDOM.render(<App />, div);
    ReactDOM.unmountComponentAtNode(div);
  });
});
