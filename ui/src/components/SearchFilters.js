import React, { useState } from "react";

const SearchFilters = ({ onIngredientsChange }) => {
  const [ingredients, setIngredients] = useState([]);
  const [inputValue, setInputValue] = useState("");

  const commonIngredients = [
    "chicken",
    "beef",
    "pork",
    "fish",
    "rice",
    "pasta",
    "potato",
    "tomato",
    "onion",
    "garlic",
    "carrot",
    "bell pepper",
    "cheese",
    "egg",
    "milk",
    "butter",
    "olive oil",
    "flour",
    "sugar",
    "salt",
  ];

  const handleAddIngredient = () => {
    if (inputValue.trim() && !ingredients.includes(inputValue.trim().toLowerCase())) {
      const newIngredients = [...ingredients, inputValue.trim().toLowerCase()];
      setIngredients(newIngredients);
      onIngredientsChange(newIngredients);
      setInputValue("");
    }
  };

  const handleRemoveIngredient = ingredient => {
    const newIngredients = ingredients.filter(item => item !== ingredient);
    setIngredients(newIngredients);
    onIngredientsChange(newIngredients);
  };

  const handleQuickAdd = ingredient => {
    if (!ingredients.includes(ingredient)) {
      const newIngredients = [...ingredients, ingredient];
      setIngredients(newIngredients);
      onIngredientsChange(newIngredients);
    }
  };

  const handleKeyDown = e => {
    if (e.key === "Enter") {
      e.preventDefault();
      handleAddIngredient();
    }
  };

  return (
    <div className="search-filters">
      <div className="ingredient-input">
        <input
          type="text"
          value={inputValue}
          onChange={e => setInputValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Add an ingredient you have"
          className="ingredient-input-field"
        />
        <button onClick={handleAddIngredient} className="add-ingredient-btn">
          Add
        </button>
      </div>

      <div className="quick-add-ingredients">
        <p>Quick add:</p>
        <div className="quick-add-buttons">
          {commonIngredients.slice(0, 10).map(ingredient => (
            <button
              key={ingredient}
              onClick={() => handleQuickAdd(ingredient)}
              className={`quick-add-btn ${ingredients.includes(ingredient) ? "selected" : ""}`}
              disabled={ingredients.includes(ingredient)}
            >
              {ingredient}
            </button>
          ))}
        </div>
      </div>

      {ingredients.length > 0 && (
        <div className="selected-ingredients">
          <p>Your ingredients:</p>
          <div className="ingredient-tags">
            {ingredients.map(ingredient => (
              <div key={ingredient} className="ingredient-tag">
                {ingredient}
                <button onClick={() => handleRemoveIngredient(ingredient)} className="remove-ingredient-btn">
                  ×
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default SearchFilters;
