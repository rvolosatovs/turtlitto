// Setup enzyme. This file will processed by create-react-app automatically.
import Enzyme from "enzyme";
import Adapter from "enzyme-adapter-react-16";
import { WebSocket } from "mock-socket";

Enzyme.configure({ adapter: new Adapter() });

global.WebSocket = WebSocket;

const realConsoleError = console.error;

beforeAll(done => {
  console.error = jest.fn().mockImplementation(msg => {
    done.fail(
      `You shouldn't log to console.error, caught the following message: ${msg}`
    );
  });
  done();
});

afterAll(() => (console.error = realConsoleError));
