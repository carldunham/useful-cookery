import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

import "./App.css";
import RecipeDetail from "./components/RecipeDetail";
import RecipeSearch from "./components/RecipeSearch";

function App() {
  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <h1>Useful Cookery</h1>
          <p>A modern, AI-enhanced cooking experience</p>
        </header>
        <main>
          <Routes>
            <Route path="/" element={<RecipeSearch userID="guest" />} />
            <Route path="/recipe/:id" element={<RecipeDetail />} />
          </Routes>
        </main>
        <footer>
          <p>&copy; {new Date().getFullYear()} Useful Cookery</p>
        </footer>
      </div>
    </Router>
  );
}

export default App;
