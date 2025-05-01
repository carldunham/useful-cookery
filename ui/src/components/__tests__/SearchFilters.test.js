import { render, screen, fireEvent } from "@testing-library/react";
import "@testing-library/jest-dom";
import userEvent from "@testing-library/user-event";
import React, { act } from "react";
import { describe, test, expect, vi, beforeEach } from "vitest";

import SearchFilters from "../SearchFilters";

describe("SearchFilters", () => {
  // Mock callback function
  const mockOnIngredientsChange = vi.fn();

  beforeEach(() => {
    // Clear mock calls between tests
    mockOnIngredientsChange.mockClear();
  });

  test("renders input field and quick add buttons", () => {
    render(<SearchFilters onIngredientsChange={mockOnIngredientsChange} />);

    // Check if input field is rendered
    expect(screen.getByPlaceholderText("Add an ingredient you have")).toBeInTheDocument();

    // Check if Add button is rendered
    expect(screen.getByText("Add")).toBeInTheDocument();

    // Check if quick add section is rendered
    expect(screen.getByText("Quick add:")).toBeInTheDocument();

    // Check if some quick add buttons are rendered
    expect(screen.getByText("chicken")).toBeInTheDocument();
    expect(screen.getByText("beef")).toBeInTheDocument();
    expect(screen.getByText("rice")).toBeInTheDocument();
  });

  test("allows adding ingredients via input field", async () => {
    render(<SearchFilters onIngredientsChange={mockOnIngredientsChange} />);

    // Get input field and Add button
    const inputField = screen.getByPlaceholderText("Add an ingredient you have");
    const addButton = screen.getByText("Add");

    // Type an ingredient and click Add
    await act(async () => {
      await userEvent.type(inputField, "broccoli");
    });

    await act(async () => {
      fireEvent.click(addButton);
    });

    // Check if ingredient was added to the list
    expect(screen.getByText("broccoli")).toBeInTheDocument();

    // Check if callback was called with correct ingredients
    expect(mockOnIngredientsChange).toHaveBeenCalledWith(["broccoli"]);

    // Input should be cleared after adding
    expect(inputField.value).toBe("");
  });

  test("allows adding ingredients by pressing Enter", async () => {
    render(<SearchFilters onIngredientsChange={mockOnIngredientsChange} />);

    // Get input field
    const inputField = screen.getByPlaceholderText("Add an ingredient you have");

    // Type an ingredient and press Enter
    await act(async () => {
      await userEvent.type(inputField, "spinach");
    });

    await act(async () => {
      fireEvent.keyDown(inputField, { key: "Enter" });
    });

    // Check if ingredient was added to the list
    expect(screen.getByText("spinach")).toBeInTheDocument();

    // Check if callback was called with correct ingredients
    expect(mockOnIngredientsChange).toHaveBeenCalledWith(["spinach"]);
  });

  test("allows adding ingredients via quick add buttons", () => {
    render(<SearchFilters onIngredientsChange={mockOnIngredientsChange} />);

    // Click on a quick add button
    fireEvent.click(screen.getByText("chicken"));

    // Check if ingredient was added to the list
    expect(screen.getByText("Your ingredients:")).toBeInTheDocument();

    // Check if the ingredient tag is in the document
    expect(screen.getByText("chicken", { selector: ".ingredient-tag" })).toBeInTheDocument();

    // Check if callback was called with correct ingredients
    expect(mockOnIngredientsChange).toHaveBeenCalledWith(["chicken"]);

    // The button should be disabled after adding
    expect(screen.getByRole("button", { name: "chicken" })).toBeDisabled();
  });

  test("allows removing ingredients", async () => {
    // Render with pre-added ingredients
    render(<SearchFilters onIngredientsChange={mockOnIngredientsChange} />);

    // Add an ingredient first
    const inputField = screen.getByPlaceholderText("Add an ingredient you have");
    await act(async () => {
      await userEvent.type(inputField, "tomato");
      fireEvent.keyDown(inputField, { key: "Enter" });
    });

    // Clear mock calls to focus on the remove action
    mockOnIngredientsChange.mockClear();

    // Find and click the remove button (×)
    const removeButton = screen.getByText("×");
    await act(async () => {
      fireEvent.click(removeButton);
    });

    // Check if ingredient tag was removed from the list
    expect(screen.queryByText("tomato", { selector: ".ingredient-tag" })).not.toBeInTheDocument();

    // Check if callback was called with empty array
    expect(mockOnIngredientsChange).toHaveBeenCalledWith([]);

    // "Your ingredients" section should not be visible when no ingredients
    expect(screen.queryByText("Your ingredients:")).not.toBeInTheDocument();
  });

  test("prevents adding duplicate ingredients", async () => {
    render(<SearchFilters onIngredientsChange={mockOnIngredientsChange} />);

    // Add an ingredient
    const inputField = screen.getByPlaceholderText("Add an ingredient you have");
    await act(async () => {
      await userEvent.type(inputField, "garlic");
      fireEvent.keyDown(inputField, { key: "Enter" });
    });

    // Clear mock calls
    mockOnIngredientsChange.mockClear();

    // Try to add the same ingredient again
    await act(async () => {
      await userEvent.type(inputField, "garlic");
      fireEvent.keyDown(inputField, { key: "Enter" });
    });

    // Check that callback wasn't called again
    expect(mockOnIngredientsChange).not.toHaveBeenCalled();

    // There should still be only one instance of the ingredient in the ingredient tags
    const garlicTags = screen.getAllByText("garlic", { selector: ".ingredient-tag" });
    expect(garlicTags.length).toBe(1);
  });

  test("trims whitespace from ingredients", async () => {
    render(<SearchFilters onIngredientsChange={mockOnIngredientsChange} />);

    // Add an ingredient with whitespace
    const inputField = screen.getByPlaceholderText("Add an ingredient you have");
    await act(async () => {
      await userEvent.type(inputField, "  lemon  ");
      fireEvent.keyDown(inputField, { key: "Enter" });
    });

    // Check if ingredient was added without whitespace
    expect(screen.getByText("lemon")).toBeInTheDocument();

    // Check if callback was called with trimmed ingredient
    expect(mockOnIngredientsChange).toHaveBeenCalledWith(["lemon"]);
  });
});
