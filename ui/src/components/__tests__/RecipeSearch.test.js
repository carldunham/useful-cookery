import { gql } from "@apollo/client";
import { MockedProvider } from "@apollo/client/testing";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import React, { act } from "react";
import "@testing-library/jest-dom";
import { describe, test, expect, vi } from "vitest";

import RecipeSearch from "../RecipeSearch";

// Mock the SearchFilters component to simplify testing
vi.mock("../SearchFilters", () => {
  return {
    default: function MockSearchFilters({ onIngredientsChange }) {
      return (
        <div data-testid="search-filters">
          <button
            onClick={() => onIngredientsChange(["tomato", "cheese"])}
            data-testid="mock-add-ingredients"
          >
            Add Ingredients
          </button>
        </div>
      );
    },
  };
});

// Mock the RecipeCard component
vi.mock("../RecipeCard", () => {
  return {
    default: function MockRecipeCard({ recipe }) {
      return <div data-testid="recipe-card">{recipe.title}</div>;
    },
  };
});

// GraphQL query definitions (must match those in the component)
const SEARCH_RECIPES = gql`
  query SearchRecipes($query: String!, $first: Int, $after: String) {
    searchRecipes(query: $query, first: $first, after: $after) {
      edges {
        node {
          id
          title
          description
          prepTime
          cookTime
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
        }
        cursor
      }
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      totalCount
    }
  }
`;

const RECOMMEND_RECIPES = gql`
  query RecommendRecipes(
    $userID: ID
    $availableIngredients: [String!]
    $first: Int
    $after: String
  ) {
    recommendRecipes(
      userID: $userID
      availableIngredients: $availableIngredients
      first: $first
      after: $after
    ) {
      edges {
        node {
          id
          title
          description
          prepTime
          cookTime
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
        }
        cursor
      }
      pageInfo {
        hasNextPage
        hasPreviousPage
        startCursor
        endCursor
      }
      totalCount
    }
  }
`;

describe("RecipeSearch", () => {
  // Mock data for search results
  const searchMocks = [
    // Mock for empty query - used when component initially renders
    {
      request: {
        query: SEARCH_RECIPES,
        variables: {
          query: "",
          first: 10,
          after: null,
        },
      },
      result: {
        data: {
          searchRecipes: {
            edges: [],
            pageInfo: {
              hasNextPage: false,
              hasPreviousPage: false,
              startCursor: null,
              endCursor: null,
            },
            totalCount: 0,
          },
        },
      },
    },
    // Original mock for "pasta" query
    {
      request: {
        query: SEARCH_RECIPES,
        variables: {
          query: "pasta",
          first: 10,
          after: null,
        },
      },
      result: {
        data: {
          searchRecipes: {
            edges: [
              {
                node: {
                  id: "recipe-1",
                  title: "Spaghetti Carbonara",
                  description: "Classic Italian pasta dish",
                  prepTime: 15,
                  cookTime: 15,
                  skillLevel: "INTERMEDIATE",
                  difficultyText: "Medium",
                  averageRating: 4.8,
                  likes: 120,
                  images: [
                    {
                      url: "carbonara.jpg",
                      alt: "Spaghetti Carbonara",
                    },
                  ],
                  author: {
                    name: "Chef Mario",
                  },
                  categories: [
                    {
                      name: "Italian",
                    },
                    {
                      name: "Pasta",
                    },
                  ],
                },
                cursor: "cursor-1",
              },
              {
                node: {
                  id: "recipe-2",
                  title: "Pasta Primavera",
                  description: "Pasta with spring vegetables",
                  prepTime: 20,
                  cookTime: 15,
                  skillLevel: "BEGINNER",
                  difficultyText: "Easy",
                  averageRating: 4.5,
                  likes: 95,
                  images: [
                    {
                      url: "primavera.jpg",
                      alt: "Pasta Primavera",
                    },
                  ],
                  author: {
                    name: "Chef Julia",
                  },
                  categories: [
                    {
                      name: "Italian",
                    },
                    {
                      name: "Vegetarian",
                    },
                  ],
                },
                cursor: "cursor-2",
              },
            ],
            pageInfo: {
              hasNextPage: true,
              hasPreviousPage: false,
              startCursor: "cursor-1",
              endCursor: "cursor-2",
            },
            totalCount: 10,
          },
        },
      },
    },
    // Original mock for "pasta" query with pagination
    {
      request: {
        query: SEARCH_RECIPES,
        variables: {
          query: "pasta",
          first: 10,
          after: "cursor-2",
        },
      },
      result: {
        data: {
          searchRecipes: {
            edges: [
              {
                node: {
                  id: "recipe-3",
                  title: "Fettuccine Alfredo",
                  description: "Creamy pasta dish",
                  prepTime: 10,
                  cookTime: 15,
                  skillLevel: "BEGINNER",
                  difficultyText: "Easy",
                  averageRating: 4.7,
                  likes: 110,
                  images: [
                    {
                      url: "alfredo.jpg",
                      alt: "Fettuccine Alfredo",
                    },
                  ],
                  author: {
                    name: "Chef Antonio",
                  },
                  categories: [
                    {
                      name: "Italian",
                    },
                    {
                      name: "Creamy",
                    },
                  ],
                },
                cursor: "cursor-3",
              },
            ],
            pageInfo: {
              hasNextPage: false,
              hasPreviousPage: true,
              startCursor: "cursor-3",
              endCursor: "cursor-3",
            },
            totalCount: 10,
          },
        },
      },
    },
    // Mock for "error" query
    {
      request: {
        query: SEARCH_RECIPES,
        variables: {
          query: "error",
          first: 10,
          after: null,
        },
      },
      error: new Error("An error occurred"),
    },
    // Mock for "nonexistent" query
    {
      request: {
        query: SEARCH_RECIPES,
        variables: {
          query: "nonexistent",
          first: 10,
          after: null,
        },
      },
      result: {
        data: {
          searchRecipes: {
            edges: [],
            pageInfo: {
              hasNextPage: false,
              hasPreviousPage: false,
              startCursor: null,
              endCursor: null,
            },
            totalCount: 0,
          },
        },
      },
    },
  ];

  // Mock data for recommendations
  const recommendMocks = [
    // Mock for empty ingredients - used when component initially renders
    {
      request: {
        query: RECOMMEND_RECIPES,
        variables: {
          userID: "test-user",
          availableIngredients: undefined,
          first: 10,
          after: null,
        },
      },
      result: {
        data: {
          recommendRecipes: {
            edges: [],
            pageInfo: {
              hasNextPage: false,
              hasPreviousPage: false,
              startCursor: null,
              endCursor: null,
            },
            totalCount: 0,
          },
        },
      },
    },
    // Original mock for recommendations with ingredients
    {
      request: {
        query: RECOMMEND_RECIPES,
        variables: {
          userID: "test-user",
          availableIngredients: ["tomato", "cheese"],
          first: 10,
          after: null,
        },
      },
      result: {
        data: {
          recommendRecipes: {
            edges: [
              {
                node: {
                  id: "recipe-4",
                  title: "Margherita Pizza",
                  description: "Classic pizza with tomato and cheese",
                  prepTime: 30,
                  cookTime: 15,
                  skillLevel: "INTERMEDIATE",
                  difficultyText: "Medium",
                  averageRating: 4.9,
                  likes: 150,
                  images: [
                    {
                      url: "pizza.jpg",
                      alt: "Margherita Pizza",
                    },
                  ],
                  author: {
                    name: "Chef Giuseppe",
                  },
                  categories: [
                    {
                      name: "Italian",
                    },
                    {
                      name: "Pizza",
                    },
                  ],
                },
                cursor: "cursor-4",
              },
              {
                node: {
                  id: "recipe-5",
                  title: "Caprese Salad",
                  description: "Simple salad with tomato and mozzarella",
                  prepTime: 10,
                  cookTime: 0,
                  skillLevel: "BEGINNER",
                  difficultyText: "Easy",
                  averageRating: 4.6,
                  likes: 85,
                  images: [
                    {
                      url: "caprese.jpg",
                      alt: "Caprese Salad",
                    },
                  ],
                  author: {
                    name: "Chef Maria",
                  },
                  categories: [
                    {
                      name: "Italian",
                    },
                    {
                      name: "Salad",
                    },
                  ],
                },
                cursor: "cursor-5",
              },
            ],
            pageInfo: {
              hasNextPage: false,
              hasPreviousPage: false,
              startCursor: "cursor-4",
              endCursor: "cursor-5",
            },
            totalCount: 2,
          },
        },
      },
    },
  ];

  test("renders search tab by default", () => {
    render(
      <MockedProvider mocks={[...searchMocks, ...recommendMocks]} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Check if search tab is active
    const searchTab = screen.getByText("Search Recipes");
    expect(searchTab).toHaveClass("active");

    // Check if search input is rendered
    expect(
      screen.getByPlaceholderText(
        "Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
      )
    ).toBeInTheDocument();
  });

  test("allows switching between search and recommend tabs", async () => {
    render(
      <MockedProvider mocks={[...searchMocks, ...recommendMocks]} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Initially on search tab
    expect(screen.getByText("Search Recipes")).toHaveClass("active");

    // Switch to recommend tab
    fireEvent.click(screen.getByText("Recommended For You"));

    // Recommend tab should be active
    expect(screen.getByText("Recommended For You")).toHaveClass("active");

    // Search input should not be visible
    expect(
      screen.queryByPlaceholderText(
        "Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
      )
    ).not.toBeInTheDocument();

    // Available Ingredients section should be visible
    expect(screen.getByText("Available Ingredients")).toBeInTheDocument();
    expect(screen.getByTestId("search-filters")).toBeInTheDocument();

    // Switch back to search tab
    fireEvent.click(screen.getByText("Search Recipes"));

    // Search tab should be active again
    expect(screen.getByText("Search Recipes")).toHaveClass("active");
  });

  test("performs search when query is entered", async () => {
    render(
      <MockedProvider mocks={searchMocks} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Enter search query
    const searchInput = screen.getByPlaceholderText(
      "Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
    );

    await act(async () => {
      await userEvent.type(searchInput, "pasta");
    });

    // Wait for search results to load
    await waitFor(() => {
      expect(screen.getByText("Spaghetti Carbonara")).toBeInTheDocument();
    });

    // Check if both recipes are displayed
    expect(screen.getByText("Spaghetti Carbonara")).toBeInTheDocument();
    expect(screen.getByText("Pasta Primavera")).toBeInTheDocument();

    // Check if recipe count is displayed
    expect(screen.getByText("Showing 2 of 10 recipes")).toBeInTheDocument();

    // Check if load more button is displayed
    expect(screen.getByText("Load More")).toBeInTheDocument();
  });

  test("loads more search results when Load More is clicked", async () => {
    render(
      <MockedProvider mocks={searchMocks} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Enter search query
    const searchInput = screen.getByPlaceholderText(
      "Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
    );

    await act(async () => {
      await userEvent.type(searchInput, "pasta");
    });

    // Wait for search results to load
    await waitFor(() => {
      expect(screen.getByText("Spaghetti Carbonara")).toBeInTheDocument();
    });

    // Click load more button
    await act(async () => {
      fireEvent.click(screen.getByText("Load More"));
    });

    // Wait for additional results to load
    await waitFor(() => {
      expect(screen.getByText("Fettuccine Alfredo")).toBeInTheDocument();
    });

    // Check if all recipes are displayed
    expect(screen.getByText("Spaghetti Carbonara")).toBeInTheDocument();
    expect(screen.getByText("Pasta Primavera")).toBeInTheDocument();
    expect(screen.getByText("Fettuccine Alfredo")).toBeInTheDocument();

    // Check if recipe count is updated
    expect(screen.getByText("Showing 3 of 10 recipes")).toBeInTheDocument();
  });

  test("displays recommendations based on ingredients", async () => {
    render(
      <MockedProvider mocks={recommendMocks} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Switch to recommend tab
    await act(async () => {
      fireEvent.click(screen.getByText("Recommended For You"));
    });

    // Add ingredients using the mocked SearchFilters component
    await act(async () => {
      fireEvent.click(screen.getByTestId("mock-add-ingredients"));
    });

    // Wait for recommendations to load
    await waitFor(() => {
      expect(screen.getByText("Margherita Pizza")).toBeInTheDocument();
    });

    // Check if both recommended recipes are displayed
    expect(screen.getByText("Margherita Pizza")).toBeInTheDocument();
    expect(screen.getByText("Caprese Salad")).toBeInTheDocument();

    // Check if recipe count is displayed
    expect(screen.getByText("Showing 2 of 2 recipes")).toBeInTheDocument();
  });

  test("displays loading state while fetching data", async () => {
    // Create a mock with a delay to ensure loading state is visible
    const delayedMock = [
      {
        request: {
          query: SEARCH_RECIPES,
          variables: {
            query: "pasta",
            first: 10,
            after: null,
          },
        },
        result: {
          data: {
            searchRecipes: {
              edges: [
                {
                  node: {
                    id: "recipe-1",
                    title: "Spaghetti Carbonara",
                    description: "Classic Italian pasta dish",
                    prepTime: 15,
                    cookTime: 15,
                    difficultyText: "Medium",
                    averageRating: 4.8,
                    likes: 120,
                    images: [
                      {
                        url: "carbonara.jpg",
                        alt: "Spaghetti Carbonara",
                      },
                    ],
                    author: {
                      name: "Chef Mario",
                    },
                    categories: [
                      {
                        name: "Italian",
                      },
                      {
                        name: "Pasta",
                      },
                    ],
                  },
                  cursor: "cursor-1",
                },
              ],
              pageInfo: {
                hasNextPage: false,
                hasPreviousPage: false,
                startCursor: "cursor-1",
                endCursor: "cursor-1",
              },
              totalCount: 1,
            },
          },
        },
        delay: 100, // Add a delay to ensure loading state is visible
      },
      // Add the empty query mock to avoid warnings
      {
        request: {
          query: SEARCH_RECIPES,
          variables: {
            query: "",
            first: 10,
            after: null,
          },
        },
        result: {
          data: {
            searchRecipes: {
              edges: [],
              pageInfo: {
                hasNextPage: false,
                hasPreviousPage: false,
                startCursor: null,
                endCursor: null,
              },
              totalCount: 0,
            },
          },
        },
      },
    ];

    render(
      <MockedProvider mocks={delayedMock} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Enter search query
    const searchInput = screen.getByPlaceholderText(
      "Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
    );

    await act(async () => {
      await userEvent.type(searchInput, "pasta");
    });

    // Loading state should be displayed
    expect(screen.getByTestId("loading-indicator")).toBeInTheDocument();

    // Wait for search results to load
    await waitFor(() => {
      expect(screen.queryByText("Loading recipes...")).not.toBeInTheDocument();
    });
  });

  test("displays error message when search fails", async () => {
    render(
      <MockedProvider mocks={searchMocks} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Enter search query that will trigger an error
    const searchInput = screen.getByPlaceholderText(
      "Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
    );

    await act(async () => {
      await userEvent.type(searchInput, "error");
    });

    // Wait for error message to appear
    await waitFor(() => {
      expect(screen.getByText("Error: An error occurred")).toBeInTheDocument();
    });
  });

  test("displays no results message when search returns empty", async () => {
    render(
      <MockedProvider mocks={searchMocks} addTypename={false}>
        <RecipeSearch userID="test-user" />
      </MockedProvider>
    );

    // Enter search query that will return no results
    const searchInput = screen.getByPlaceholderText(
      "Search recipes with natural language (e.g., 'vegetarian pasta without mushrooms')"
    );

    await act(async () => {
      await userEvent.type(searchInput, "nonexistent");
    });

    // Wait for no results message to appear
    await waitFor(() => {
      expect(
        screen.getByText("No recipes found. Try a different search query.")
      ).toBeInTheDocument();
    });
  });
});
