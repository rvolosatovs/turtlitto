// Setup enzyme. This file will processed by create-react-app automatically.
import Enzyme from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import { WebSocket } from "mock-socket";

Enzyme.configure({ adapter: new Adapter() });

global.WebSocket = WebSocket;
