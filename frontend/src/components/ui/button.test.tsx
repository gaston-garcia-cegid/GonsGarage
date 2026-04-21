import React from "react";
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Button } from "./button";

describe("Button (shadcn)", () => {
  it("invokes onClick when the user activates the default button", async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();
    render(
      <Button type="button" onClick={onClick}>
        Guardar
      </Button>,
    );
    const btn = screen.getByRole("button", { name: "Guardar" });
    expect(btn).not.toBeDisabled();
    await user.click(btn);
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it("keeps destructive actions disabled when disabled is true", async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();
    const { rerender } = render(
      <Button type="button" variant="destructive" onClick={onClick}>
        Eliminar
      </Button>,
    );
    const btn = screen.getByRole("button", { name: "Eliminar" });
    expect(btn).not.toBeDisabled();
    await user.click(btn);
    expect(onClick).toHaveBeenCalledTimes(1);

    rerender(
      <Button type="button" variant="destructive" disabled onClick={onClick}>
        Eliminar
      </Button>,
    );
    const disabledBtn = screen.getByRole("button", { name: "Eliminar" });
    expect(disabledBtn).toBeDisabled();
    await user.click(disabledBtn);
    expect(onClick).toHaveBeenCalledTimes(1);
  });
});
