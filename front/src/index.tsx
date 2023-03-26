import React from 'react';
import ReactDOM from 'react-dom/client';
import Header from './header';
import Mtr from './mtr';
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
