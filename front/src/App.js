import React, { Component } from "react";
import "./App.css";
import "normalize.css";
import styled from "styled-components";
import SRRButton from "./SRRButton";

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

const SpecialP = styled.p`
  margin: 0px;
  padding: 0px;
`;

class App extends Component {
  render() {
    return (
      <AppWrap>
        <Footer>
          <StartButton
            buttonText={<SpecialP>&#9658;</SpecialP>}
            onClick={() => {
              console.log("hurrrr");
            }}
            enabled={true}
          />
          <SettingsButton
            buttonText={<SpecialP>Settings</SpecialP>}
            onClick={() => {
              console.log("harrrr");
            }}
            enabled={true}
          />
          <StopButton
            buttonText={<SpecialP>&#9724;</SpecialP>}
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
