export default function ThemeToggle() {
  const [isDark, setIsDark] = useState(false);

  useEffect(() => {
    const savedTheme = localStorage.getItem("theme");
    const systemDark = window.matchMedia(
      "(prefers-color-scheme: dark)",
    ).matches;

    if (savedTheme === "dark" || (!savedTheme && systemDark)) {
      setIsDark(true);
      document.documentElement.classList.add("dark");
    }
  }, []);

  const toggleTheme = () => {
    setIsDark(!isDark);
    if (!isDark) {
      document.documentElement.classList.add("dark");
      localStorage.setItem("theme", "dark");
    } else {
      document.documentElement.classList.remove("dark");
      localStorage.setItem("theme", "light");
    }
  };

  return (
    <button
      onClick={toggleTheme}
      title={isDark ? "Switch to Light Mode" : "Switch to Dark Mode"}
      className="flex items-center gap-2 px-2 py-1 rounded-full transition-colors"
    >
      {isDark ? (
        <>
          <Moon size={16} className="text-indigo-400" />
          <span className="text-xs font-medium text-[var(--text-primary)]">
            Dark
          </span>
        </>
      ) : (
        <>
          <Sun size={16} className="text-amber-500" />
          <span className="text-xs font-medium text-[var(--text-primary)]">
            Light
          </span>
        </>
      )}
    </button>
  );
}

import { Sun, Moon } from "lucide-react";
import { useEffect, useState } from "react";
