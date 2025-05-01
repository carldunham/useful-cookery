import React from "react";
import { useNavigate } from "react-router-dom";
import "./RecipeCard.css";

const RecipeCard = ({ recipe }) => {
  const navigate = useNavigate();
  const {
    title,
    description,
    prepTime,
    cookTime,
    difficultyText,
    skillLevel,
    averageRating,
    images,
    author,
    categories,
  } = recipe;

  // Default image if none provided
  const imageUrl =
    images && images.length > 0 && images[0].url
      ? images[0].url
      : "https://via.placeholder.com/300x200?text=No+Image";
  const imageAlt = images && images.length > 0 ? images[0].alt : "Recipe image";

  // Format categories
  const categoryNames = categories ? categories.map((cat) => cat.name).join(", ") : "";

  // Calculate total time
  const totalTime = prepTime && cookTime ? prepTime + cookTime : null;

  return (
    <div className="recipe-card">
      <div className="recipe-image">
        <img src={imageUrl} alt={imageAlt} />
      </div>
      <div className="recipe-content">
        <h3 className="recipe-title">{title}</h3>
        {author && <p className="recipe-author">By {author.name}</p>}
        <p className="recipe-description">{description}</p>
        <div className="recipe-meta">
          {skillLevel && <span className="recipe-skill-level">Skill Level: {skillLevel}</span>}
          {difficultyText && (
            <span className="recipe-difficulty">Difficulty: {difficultyText}</span>
          )}
          {totalTime && <span className="recipe-time">Time: {totalTime} min</span>}
          {averageRating && (
            <span className="recipe-rating">Rating: {averageRating.toFixed(1)}/5</span>
          )}
        </div>
        {categoryNames && <p className="recipe-categories">Categories: {categoryNames}</p>}
        <button
          className="view-recipe-btn"
          onClick={() => navigate(`/recipe/${recipe.id}`)}
        >
          View Recipe
        </button>
      </div>
    </div>
  );
};

export default RecipeCard;
