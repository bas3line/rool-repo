const root = document.documentElement;
const themeColor = document.querySelector('meta[name="theme-color"]');
const themeToggles = [...document.querySelectorAll("[data-theme-toggle]")];

const syncTheme = () => {
  const isDark = root.dataset.theme === "dark";
  themeColor?.setAttribute("content", isDark ? "#111111" : "#fafafa");
  themeToggles.forEach((toggle) => {
    toggle.setAttribute("aria-label", `Switch to ${isDark ? "light" : "dark"} mode`);
  });
};

themeToggles.forEach((toggle) => {
  toggle.addEventListener("click", () => {
    const nextTheme = root.dataset.theme === "dark" ? "light" : "dark";
    root.dataset.theme = nextTheme;
    try {
      localStorage.setItem("portfolio-theme-v3", nextTheme);
    } catch {}
    syncTheme();
  });
});

syncTheme();

document.querySelectorAll("[data-copy]").forEach((button) => {
  button.addEventListener("click", async () => {
    const value = button.dataset.copy;
    try {
      await navigator.clipboard.writeText(value);
    } catch {
      const textarea = document.createElement("textarea");
      textarea.value = value;
      textarea.setAttribute("readonly", "");
      textarea.style.position = "fixed";
      textarea.style.opacity = "0";
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand("copy");
      textarea.remove();
    }

    button.textContent = "copied";
    button.dataset.copied = "true";
    window.setTimeout(() => {
      button.textContent = "copy";
      delete button.dataset.copied;
    }, 1500);
  });
});
