import * as React from 'react';
import * as ReactDOM from 'react-dom';
import {App} from "./App";
import * as Wails from "@wailsapp/runtime";

import './styles.css';

Wails.Init(() => {
    ReactDOM.render(
        <App/>,
        document.getElementById('app')
    )
});