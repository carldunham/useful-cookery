import { renderHook, act } from "@testing-library/react";
import { describe, test, expect, vi, beforeEach, afterEach } from "vitest";

import { useDebounce } from "../useDebounce";

// Mock timer functions
beforeEach(() => {
  vi.useFakeTimers();
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("useDebounce", () => {
  test("should return the initial value immediately", () => {
    const { result } = renderHook(() => useDebounce("initial value", 500));

    expect(result.current).toBe("initial value");
  });

  test("should not update the value before the delay has passed", () => {
    const { result, rerender } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: "initial value", delay: 500 },
    });

    // Change the value
    rerender({ value: "updated value", delay: 500 });

    // Value should not have changed yet
    expect(result.current).toBe("initial value");

    // Fast-forward time by 400ms (less than the delay)
    act(() => {
      vi.advanceTimersByTime(400);
    });

    // Value should still not have changed
    expect(result.current).toBe("initial value");
  });

  test("should update the value after the delay has passed", () => {
    const { result, rerender } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: "initial value", delay: 500 },
    });

    // Change the value
    rerender({ value: "updated value", delay: 500 });

    // Fast-forward time by 500ms (equal to the delay)
    act(() => {
      vi.advanceTimersByTime(500);
    });

    // Value should have updated
    expect(result.current).toBe("updated value");
  });

  test("should handle multiple value changes within the delay period", () => {
    const { result, rerender } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: "initial value", delay: 500 },
    });

    // Change the value multiple times
    rerender({ value: "intermediate value 1", delay: 500 });

    // Fast-forward time by 200ms
    act(() => {
      vi.advanceTimersByTime(200);
    });

    // Change the value again
    rerender({ value: "intermediate value 2", delay: 500 });

    // Fast-forward time by 200ms
    act(() => {
      vi.advanceTimersByTime(200);
    });

    // Change the value one more time
    rerender({ value: "final value", delay: 500 });

    // Value should still be the initial value
    expect(result.current).toBe("initial value");

    // Fast-forward time by 500ms to complete the delay for the last change
    act(() => {
      vi.advanceTimersByTime(500);
    });

    // Value should be the final value
    expect(result.current).toBe("final value");
  });

  test("should handle delay changes", () => {
    const { result, rerender } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: "initial value", delay: 500 },
    });

    // Change the value and delay
    rerender({ value: "updated value", delay: 1000 });

    // Fast-forward time by 500ms (the original delay)
    act(() => {
      vi.advanceTimersByTime(500);
    });

    // Value should not have changed yet due to increased delay
    expect(result.current).toBe("initial value");

    // Fast-forward time by another 500ms to complete the new delay
    act(() => {
      vi.advanceTimersByTime(500);
    });

    // Value should have updated
    expect(result.current).toBe("updated value");
  });

  test("should clean up timeout on unmount", () => {
    const clearTimeoutSpy = vi.spyOn(global, "clearTimeout");

    const { unmount } = renderHook(({ value, delay }) => useDebounce(value, delay), {
      initialProps: { value: "test value", delay: 500 },
    });

    // Unmount the component
    unmount();

    // clearTimeout should have been called
    expect(clearTimeoutSpy).toHaveBeenCalled();
  });
});
