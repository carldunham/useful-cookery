import { render, screen } from "@testing-library/react";
import React from "react";
import { describe, test, expect } from "vitest";

import "@testing-library/jest-dom";
import RecipeCard from "../RecipeCard";

describe("RecipeCard", () => {
  const mockRecipe = {
    id: "recipe-1",
    title: "Test Recipe",
    description: "This is a test recipe description",
    prepTime: 15,
    cookTime: 30,
    difficulty: "Medium",
    averageRating: 4.5,
    images: [{ url: "test-image.jpg", alt: "Test recipe image" }],
    author: { name: "Test Author" },
    categories: [{ name: "Italian" }, { name: "Pasta" }],
  };

  test("renders recipe information correctly", () => {
    render(<RecipeCard recipe={mockRecipe} />);

    // Check if title is rendered
    expect(screen.getByText("Test Recipe")).toBeInTheDocument();

    // Check if description is rendered
    expect(screen.getByText("This is a test recipe description")).toBeInTheDocument();

    // Check if author is rendered
    expect(screen.getByText("By Test Author")).toBeInTheDocument();

    // Check if difficulty is rendered
    expect(screen.getByText("Difficulty: Medium")).toBeInTheDocument();

    // Check if total time is rendered
    expect(screen.getByText("Time: 45 min")).toBeInTheDocument();

    // Check if rating is rendered
    expect(screen.getByText("Rating: 4.5/5")).toBeInTheDocument();

    // Check if categories are rendered
    expect(screen.getByText("Categories: Italian, Pasta")).toBeInTheDocument();

    // Check if image is rendered with correct attributes
    const image = screen.getByAltText("Test recipe image");
    expect(image).toBeInTheDocument();
    expect(image).toHaveAttribute("src", "test-image.jpg");

    // Check if view recipe button is rendered
    expect(screen.getByText("View Recipe")).toBeInTheDocument();
  });

  test("handles missing data gracefully", () => {
    const incompleteRecipe = {
      title: "Incomplete Recipe",
      description: "This recipe has missing data",
    };

    render(<RecipeCard recipe={incompleteRecipe} />);

    // Check if title is rendered
    expect(screen.getByText("Incomplete Recipe")).toBeInTheDocument();

    // Check if description is rendered
    expect(screen.getByText("This recipe has missing data")).toBeInTheDocument();

    // Check if default image is used
    const image = screen.getByAltText("Recipe image");
    expect(image).toBeInTheDocument();
    expect(image).toHaveAttribute("src", "https://via.placeholder.com/300x200?text=No+Image");

    // Check that optional elements are not rendered
    expect(screen.queryByText(/By/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Difficulty/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Time/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Rating/)).not.toBeInTheDocument();
    expect(screen.queryByText(/Categories/)).not.toBeInTheDocument();
  });
});
