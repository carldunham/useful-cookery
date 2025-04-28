import React, { useState } from "react";
import RecipeCard from "../components/RecipeCard";
import SearchFilters from "../components/SearchFilters";
import { useDebounce } from "../hooks/useDebounce";
// Import the generated hooks
// Note: These hooks will be generated with the prefix "Example" once the code generator is run
import { useExampleSearchRecipesLazyQuery, useExampleRecommendRecipesQuery } from "../generated/graphql";

const RecipeSearchWithGeneratedHooks = ({ userID }) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [availableIngredients, setAvailableIngredients] = useState([]);
  const [activeTab, setActiveTab] = useState("search"); // 'search' or 'recommend'

  // Debounce search to prevent too many API calls
  const debouncedSearchQuery = useDebounce(searchQuery, 500);

  // Pagination state
  const [pageSize] = useState(10);
  const [searchAfter, setSearchAfter] = useState(null);
  const [recommendAfter, setRecommendAfter] = useState(null);

  // Lazy query for search using generated hook
  const [executeSearch, { loading: searchLoading, error: searchError, data: searchData }] =
    useExampleSearchRecipesLazyQuery({
      variables: {
        query: debouncedSearchQuery,
        first: pageSize,
        after: searchAfter,
      },
      skip: !debouncedSearchQuery,
    });

  // Get recommendations using generated hook
  const {
    loading: recommendLoading,
    error: recommendError,
    data: recommendData,
  } = useExampleRecommendRecipesQuery({
    variables: {
      userID,
      availableIngredients: availableIngredients.length > 0 ? availableIngredients : undefined,
      first: pageSize,
      after: recommendAfter,
    },
    skip: activeTab !== "recommend",
  });

  // Handle search input change
  const handleSearchChange = e => {
    const newQuery = e.target.value;
    setSearchQuery(newQuery);

    if (newQuery.length >= 3) {
      executeSearch();
    }
  };

  // Handle ingredient changes
  const handleIngredientChange = ingredients => {
    setAvailableIngredients(ingredients);
  };

  // Handle tab change
  const handleTabChange = tab => {
    setActiveTab(tab);
  };

  // Determine which recipes to display
  const getRecipesToDisplay = () => {
    if (activeTab === "search") {
      return searchData?.searchRecipes?.edges?.map(edge => edge.node) || [];
    } else {
      return recommendData?.recommendRecipes?.edges?.map(edge => edge.node) || [];
    }
  };

  // Get current page info
  const getCurrentPageInfo = () => {
    if (activeTab === "search") {
      return searchData?.searchRecipes?.pageInfo;
    } else {
      return recommendData?.recommendRecipes?.pageInfo;
    }
  };

  // Get total count
  const getTotalCount = () => {
    if (activeTab === "search") {
      return searchData?.searchRecipes?.totalCount || 0;
    } else {
      return recommendData?.recommendRecipes?.totalCount || 0;
    }
  };

  // Handle load more
  const handleLoadMore = () => {
    const pageInfo = getCurrentPageInfo();
    if (pageInfo?.hasNextPage) {
      if (activeTab === "search") {
        setSearchAfter(pageInfo.endCursor);
        executeSearch({
          variables: {
            query: debouncedSearchQuery,
            first: pageSize,
            after: pageInfo.endCursor,
          },
        });
      } else {
        setRecommendAfter(pageInfo.endCursor);
      }
    }
  };

  // Check if loading
  const isLoading = activeTab === "search" ? searchLoading : recommendLoading;

  // Check for errors
  const error = activeTab === "search" ? searchError : recommendError;

  return (
    <div className="recipe-search">
      <div className="search-tabs">
        <button className={`tab ${activeTab === "search" ? "active" : ""}`} onClick={() => handleTabChange("search")}>
          Search Recipes
        </button>
        <button
          className={`tab ${activeTab === "recommend" ? "active" : ""}`}
          onClick={() => handleTabChange("recommend")}
        >
          Recommended For You
        </button>
      </div>

      {activeTab === "search" && (
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

      {activeTab === "recommend" && (
        <div className="ingredient-filter">
          <h3>Available Ingredients</h3>
          <SearchFilters onIngredientsChange={handleIngredientChange} />
        </div>
      )}

      {isLoading && <div className="loading">Loading recipes...</div>}

      {error && <div className="error">Error: {error.message}</div>}

      {!isLoading && !error && (
        <div className="recipe-results">
          {getRecipesToDisplay().length === 0 ? (
            <div className="no-results">
              {activeTab === "search"
                ? "No recipes found. Try a different search query."
                : "No recommendations available. Try adding some ingredients."}
            </div>
          ) : (
            <>
              <div className="recipe-count">
                Showing {getRecipesToDisplay().length} of {getTotalCount()} recipes
              </div>
              <div className="recipe-grid">
                {getRecipesToDisplay().map(recipe => (
                  <RecipeCard key={recipe.id} recipe={recipe} />
                ))}
              </div>
              {getCurrentPageInfo()?.hasNextPage && (
                <div className="load-more">
                  <button onClick={handleLoadMore} disabled={isLoading}>
                    {isLoading ? "Loading..." : "Load More"}
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default RecipeSearchWithGeneratedHooks;
