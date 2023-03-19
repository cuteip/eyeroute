import React from 'react';
import ReactDOM from 'react-dom/client';
import Header from './header';
import Mtr from './mtr';
import reportWebVitals from './reportWebVitals';
import { Routes, Route, BrowserRouter } from "react-router-dom";

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <div>
    <BrowserRouter>
      <Header />
      <Routes>
        <Route index element={<div></div>} />
        <Route path="mtr" element={<Mtr />} />
      </Routes>
    </BrowserRouter>
  </div>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
