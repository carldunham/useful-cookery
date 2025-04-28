import React from "react";

import "./App.css";
import RecipeSearch from "./components/RecipeSearch";

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Useful Cookery</h1>
        <p>A modern, AI-enhanced cooking experience</p>
      </header>
      <main>
        <RecipeSearch userID="guest" />
      </main>
      <footer>
        <p>&copy; {new Date().getFullYear()} Useful Cookery</p>
      </footer>
    </div>
  );
}

export default App;
