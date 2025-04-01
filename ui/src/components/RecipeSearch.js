import React, { useState } from 'react';
import { useQuery, useLazyQuery } from '@apollo/client';
import { gql } from '@apollo/client';
import RecipeCard from './RecipeCard';
import SearchFilters from './SearchFilters';
import { useDebounce } from '../hooks/useDebounce';

// GraphQL query definitions
const SEARCH_RECIPES = gql`
  query SearchRecipes($query: String!) {
    searchRecipes(query: $query) {
      id
      title
      description
      prepTime
      cookTime
      difficulty
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
    }
  }
`;

const RECOMMEND_RECIPES = gql`
  query RecommendRecipes($userId: ID, $availableIngredients: [String!]) {
    recommendRecipes(userId: $userId, availableIngredients: $availableIngredients) {
      id
      title
      description
      prepTime
      cookTime
      difficulty
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
    }
  }
`;

const RecipeSearch = ({ userId }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [availableIngredients, setAvailableIngredients] = useState([]);
  const [activeTab, setActiveTab] = useState('search'); // 'search' or 'recommend'
  
  // Debounce search to prevent too many API calls
  const debouncedSearchQuery = useDebounce(searchQuery, 500);
  
  // Lazy query for search
  const [executeSearch, { loading: searchLoading, error: searchError, data: searchData }] = useLazyQuery(
    SEARCH_RECIPES,
    {
      variables: { query: debouncedSearchQuery },
      skip: !debouncedSearchQuery
    }
  );
  
  // Get recommendations
  const { loading: recommendLoading, error: recommendError, data: recommendData } = useQuery(
    RECOMMEND_RECIPES,
    {
      variables: { 
        userId, 
        availableIngredients: availableIngredients.length > 0 ? availableIngredients : undefined 
      },
      skip: activeTab !== 'recommend'
    }
  );
  
  // Handle search input change
  const handleSearchChange = (e) => {
    const newQuery = e.target.value;
    setSearchQuery(newQuery);
    
    if (newQuery.length >= 3) {
      executeSearch();
    }
  };
  
  // Handle ingredient changes
  const handleIngredientChange = (ingredients) => {
    setAvailableIngredients(ingredients);
  };
  
  // Handle tab change
  const handleTabChange = (tab) => {
    setActiveTab(tab);
  };
  
  // Determine which recipes to display
  const getRecipesToDisplay = () => {
    if (activeTab === 'search') {
      return searchData?.searchRecipes || [];
    } else {
      return recommendData?.recommendRecipes || [];
    }
  };
  
  // Check if loading
  const isLoading = activeTab === 'search' ? searchLoading : recommendLoading;
  
  // Check for errors
  const error = activeTab === 'search' ? searchError : recommendError;
  
  return (
    <div className="recipe-search">
      <div className="search-tabs">
        <button
          className={`tab ${activeTab === 'search' ? 'active' : ''}`}
          onClick={() => handleTabChange('search')}
        >
          Search Recipes
        </button>
        <button
          className={`tab ${activeTab === 'recommend' ? 'active' : ''}`}
          onClick={() => handleTabChange('recommend')}
        >
          Recommended For You
        </button>
      </div>
      
      {activeTab === 'search' && (
        <div className="search-input-container">
          <input
            type="text"
            placeholder="Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
            value={searchQuery}
            onChange={handleSearchChange}
            className="search-input"
          />
        </div>
      )}
      
      {activeTab === 'recommend' && (
        <div className="ingredient-filter">
          <h3>Available Ingredients</h3>
          <SearchFilters onIngredientsChange={handleIngredientChange} />
        </div>
      )}
      
      {isLoading && <div className="loading">Loading recipes...</div>}
      
      {error && (
        <div className="error">
          Error: {error.message}
        </div>
      )}
      
      {!isLoading && !error && (
        <div className="recipe-results">
          {getRecipesToDisplay().length === 0 ? (
            <div className="no-results">
              {activeTab === 'search' 
                ? 'No recipes found. Try a different search query.'
                : 'No recommendations available. Try adding some ingredients.'}
            </div>
          ) : (
            <div className="recipe-grid">
              {getRecipesToDisplay().map(recipe => (
                <RecipeCard key={recipe.id} recipe={recipe} />
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default RecipeSearch;
