import React, { Component } from 'react';
import styled from "styled-components";
import TurtleEnableButton from './TurtleEnableButton';

/**
 * Simple bar with all the turtles
 * Author: S.A. Tanja
 */
const TurtleEnableBar = () => {
    return(
        <div>
            <TurtleEnableButton id = "1"/>
            <TurtleEnableButton id = "2"/>
            <TurtleEnableButton id = "3"/>
            <TurtleEnableButton id = "4"/>
            <TurtleEnableButton id = "5"/>
            <TurtleEnableButton id = "6"/>
        </div>
    );
}

export default TurtleEnableBar;