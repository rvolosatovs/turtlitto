import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import styled from "styled-components";
import SRRButton from "./SRRButton";
import Turtle from "./Turtle";

const AppWrap = styled.div`\
  width = 100%;
  height = 100%;
`;

const Footer = styled.footer`
  position: fixed;
  bottom: 0;
  width: 100%;
  min-width: 100px;
  height: 10%;
  border-style: solid;
  border-width: 5px 0px 0px 0px;
  display: flex;
  justify-content: space-between;
  margin: 0px;
`;

const StartButton = styled(SRRButton)`
  width: 32%;
  font-size: 4vw;
  content: "\&#9658;";
`;

const SettingsButton = styled(SRRButton)`
  width: 32%;
  font-size: 3vw;
`;

const StopButton = styled(SRRButton)`
  width: 32%;
  font-size: 4vw;
`;

const AugmentedText = styled.p`
  margin: 0px;
  padding: 0px;
`;

const SendCommand = command => {};

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      //messages: [],
      isConnected: false,
      port: "4242",
      host: "localhost"
    };

    this.connection = null;
  }
  render() {
    return (
      <AppWrap>
        <Turtle
          turtle={{
            battery: 100,
            editable: true,
            role: "Goalkeeper",
            home: "Yellow home",
            team: "Magenta",
            id: 2
          }}
        />

        <Footer>
          <StartButton
            buttonText={<AugmentedText>&#9658;</AugmentedText>}
            onClick={() => {
              console.log("hurrrr");
            }}
            enabled={true}
          />
          <SettingsButton
            buttonText={<AugmentedText>Settings</AugmentedText>}
            onClick={() => {
              console.log("harrrr");
            }}
            enabled={true}
          />
          <StopButton
            buttonText={<AugmentedText>&#9724;</AugmentedText>}
            onClick={() => {
              console.log("hrrrrr");
            }}
            enabled={true}
          />
        </Footer>
      </AppWrap>
    );
  }
}

export default App;
