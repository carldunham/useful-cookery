import { MockedProvider } from "@apollo/client/testing";
import { render, screen } from "@testing-library/react";
import React, { act } from "react";
import "@testing-library/jest-dom";

import App from "../App";

// Mock the RecipeSearch component to simplify testing
jest.mock("../components/RecipeSearch", () => {
  return function MockRecipeSearch({ userID }) {
    return <div data-testid="recipe-search">Recipe Search (User ID: {userID})</div>;
  };
});

describe("App", () => {
  test("renders header with title", () => {
    render(
      <MockedProvider>
        <App />
      </MockedProvider>
    );

    // Check if the header title is rendered
    expect(screen.getByText("Useful Cookery")).toBeInTheDocument();

    // Check if the tagline is rendered
    expect(screen.getByText("A modern, AI-enhanced cooking experience")).toBeInTheDocument();
  });

  test("renders RecipeSearch component with guest user ID", () => {
    render(
      <MockedProvider>
        <App />
      </MockedProvider>
    );

    // Check if RecipeSearch is rendered with the correct props
    expect(screen.getByTestId("recipe-search")).toBeInTheDocument();
    expect(screen.getByText("Recipe Search (User ID: guest)")).toBeInTheDocument();
  });

  test("renders footer with copyright information", () => {
    // Mock the Date object to return a fixed year
    const originalDate = global.Date;
    global.Date = class extends Date {
      getFullYear() {
        return 2025;
      }
    };

    render(
      <MockedProvider>
        <App />
      </MockedProvider>
    );

    // Check if footer with copyright is rendered
    expect(screen.getByText("© 2025 Useful Cookery")).toBeInTheDocument();

    // Restore the original Date object
    global.Date = originalDate;
  });
});
