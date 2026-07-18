const copyButtons = document.querySelectorAll("[data-copy]");

copyButtons.forEach((button) => {
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

const tocLinks = [...document.querySelectorAll(".toc a")];
const sections = tocLinks
  .map((link) => document.querySelector(link.getAttribute("href")))
  .filter(Boolean);

if ("IntersectionObserver" in window) {
  const observer = new IntersectionObserver((entries) => {
    const visible = entries
      .filter((entry) => entry.isIntersecting)
      .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)[0];
    if (!visible) return;

    tocLinks.forEach((link) => {
      const active = link.getAttribute("href") === `#${visible.target.id}`;
      link.toggleAttribute("aria-current", active);
    });
  }, { rootMargin: "-15% 0px -70% 0px" });

  sections.forEach((section) => observer.observe(section));
}
