import { useQuery } from "@apollo/client";
import { gql } from "@apollo/client";
import React from "react";
import { useParams, Link } from "react-router-dom";
import "./RecipeDetail.css";

// GraphQL query to fetch a single recipe by ID
const GET_RECIPE = gql`
  query GetRecipe($id: ID!) {
    recipe(id: $id) {
      id
      title
      description
      prepTime
      cookTime
      servings
      difficultyText
      skillLevel
      averageRating
      likes
      images {
        url
        alt
      }
      author {
        name
      }
      categories {
        name
      }
      ingredients {
        name
        units {
          value
          unit
          system
        }
        preparation
        isOptional
        substitutes
      }
      steps {
        orderIndex
        description
        timeEstimate
        image {
          url
          alt
        }
      }
      nutritionInfo {
        calories
        protein
        carbs
        fat
        fiber
        sugar
        sodium
      }
      tags
      createdAt
      updatedAt
    }
  }
`;

const RecipeDetail = () => {
  // Get the recipe ID from the URL
  const { id } = useParams();

  // Fetch the recipe data
  const { loading, error, data } = useQuery(GET_RECIPE, {
    variables: { id },
    onError: (error) => {
      console.error("Error fetching recipe:", error);
    },
    onCompleted: (data) => {
      console.log("Recipe data received:", data);
    },
  });

  // Enhanced error handling with more details
  if (loading) return <div className="loading">Loading recipe...</div>;
  if (error) {
    console.error("Detailed error:", error);
    return (
      <div className="error">
        <h2>Error Loading Recipe</h2>
        <p>{error.message}</p>
        {error.graphQLErrors?.map((err, i) => (
          <p key={i}>GraphQL Error: {err.message}</p>
        ))}
        {error.networkError && <p>Network Error: {error.networkError.message}</p>}
        <Link to="/" className="back-link">
          &larr; Back to Recipes
        </Link>
      </div>
    );
  }
  if (!data || !data.recipe) {
    console.log("No recipe data found for ID:", id);
    return (
      <div className="not-found">
        <h2>Recipe Not Found</h2>
        <p>We couldn't find the recipe you're looking for.</p>
        <Link to="/" className="back-link">
          &larr; Back to Recipes
        </Link>
      </div>
    );
  }

  const recipe = data.recipe;

  // Calculate total time
  const totalTime = recipe.prepTime && recipe.cookTime ? recipe.prepTime + recipe.cookTime : null;

  // Format categories
  const categoryNames = recipe.categories ? recipe.categories.map((cat) => cat.name).join(", ") : "";

  // Default image if none provided
  const mainImage =
    recipe.images && recipe.images.length > 0
      ? recipe.images[0]
      : { url: "https://via.placeholder.com/600x400?text=No+Image", alt: "Recipe image" };

  return (
    <div className="recipe-detail">
      <Link to="/" className="back-link">
        &larr; Back to Recipes
      </Link>

      <div className="recipe-header">
        <h1 className="recipe-title">{recipe.title}</h1>
        {recipe.author && <p className="recipe-author">By {recipe.author.name}</p>}
      </div>

      <div className="recipe-image-container">
        <img src={mainImage.url} alt={mainImage.alt} className="recipe-main-image" />
      </div>

      <div className="recipe-meta">
        {recipe.skillLevel && (
          <span className="recipe-skill-level">Skill Level: {recipe.skillLevel}</span>
        )}
        {recipe.difficultyText && (
          <span className="recipe-difficulty">Difficulty: {recipe.difficultyText}</span>
        )}
        {totalTime && <span className="recipe-time">Total Time: {totalTime} min</span>}
        {recipe.prepTime && <span className="recipe-prep-time">Prep: {recipe.prepTime} min</span>}
        {recipe.cookTime && <span className="recipe-cook-time">Cook: {recipe.cookTime} min</span>}
        {recipe.servings && <span className="recipe-servings">Servings: {recipe.servings}</span>}
        {recipe.averageRating && (
          <span className="recipe-rating">Rating: {recipe.averageRating.toFixed(1)}/5</span>
        )}
        {recipe.likes && <span className="recipe-likes">Likes: {recipe.likes}</span>}
      </div>

      {categoryNames && <p className="recipe-categories">Categories: {categoryNames}</p>}

      <div className="recipe-description">
        <h2>Description</h2>
        <p>{recipe.description}</p>
      </div>

      <div className="recipe-ingredients">
        <h2>Ingredients</h2>
        <ul>
          {recipe.ingredients.map((ingredient, index) => {
            // Get the main unit (or first if no main unit is specified)
            const mainUnit = ingredient.units?.find((u) => u.isMain) || ingredient.units?.[0];

            return (
              <li key={index} className={ingredient.isOptional ? "optional" : ""}>
                {mainUnit && (
                  <span className="ingredient-quantity">
                    {mainUnit.value} {mainUnit.unit}
                  </span>
                )}{" "}
                <span className="ingredient-name">{ingredient.name}</span>
                {ingredient.preparation && (
                  <span className="ingredient-prep">, {ingredient.preparation}</span>
                )}
                {ingredient.isOptional && <span className="optional-label"> (optional)</span>}
              </li>
            );
          })}
        </ul>
      </div>

      <div className="recipe-steps">
        <h2>Instructions</h2>
        <ol>
          {[...recipe.steps]
            .sort((a, b) => a.orderIndex - b.orderIndex)
            .map((step, index) => (
              <li key={index}>
                <div className="step-content">
                  <p>{step.description}</p>
                  {step.timeEstimate && (
                    <span className="step-time">
                      Approximately {step.timeEstimate} minutes
                    </span>
                  )}
                </div>
                {step.image && (
                  <div className="step-image">
                    <img src={step.image.url} alt={step.image.alt || `Step ${index + 1}`} />
                  </div>
                )}
              </li>
            ))}
        </ol>
      </div>

      {recipe.nutritionInfo && (
        <div className="recipe-nutrition">
          <h2>Nutrition Information</h2>
          <div className="nutrition-grid">
            {recipe.nutritionInfo.calories && (
              <div className="nutrition-item">
                <span className="nutrition-label">Calories</span>
                <span className="nutrition-value">{recipe.nutritionInfo.calories}</span>
              </div>
            )}
            {recipe.nutritionInfo.protein && (
              <div className="nutrition-item">
                <span className="nutrition-label">Protein</span>
                <span className="nutrition-value">{recipe.nutritionInfo.protein}g</span>
              </div>
            )}
            {recipe.nutritionInfo.carbs && (
              <div className="nutrition-item">
                <span className="nutrition-label">Carbs</span>
                <span className="nutrition-value">{recipe.nutritionInfo.carbs}g</span>
              </div>
            )}
            {recipe.nutritionInfo.fat && (
              <div className="nutrition-item">
                <span className="nutrition-label">Fat</span>
                <span className="nutrition-value">{recipe.nutritionInfo.fat}g</span>
              </div>
            )}
            {recipe.nutritionInfo.fiber && (
              <div className="nutrition-item">
                <span className="nutrition-label">Fiber</span>
                <span className="nutrition-value">{recipe.nutritionInfo.fiber}g</span>
              </div>
            )}
            {recipe.nutritionInfo.sugar && (
              <div className="nutrition-item">
                <span className="nutrition-label">Sugar</span>
                <span className="nutrition-value">{recipe.nutritionInfo.sugar}g</span>
              </div>
            )}
            {recipe.nutritionInfo.sodium && (
              <div className="nutrition-item">
                <span className="nutrition-label">Sodium</span>
                <span className="nutrition-value">{recipe.nutritionInfo.sodium}mg</span>
              </div>
            )}
          </div>
        </div>
      )}

      {recipe.tags && recipe.tags.length > 0 && (
        <div className="recipe-tags">
          <h2>Tags</h2>
          <div className="tags-container">
            {recipe.tags.map((tag, index) => (
              <span key={index} className="tag">
                {tag}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default RecipeDetail;
