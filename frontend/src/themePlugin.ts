import {type App, inject, ref, type Ref} from 'vue';

type Theme = string;

interface ThemePluginOptions {
    defaultTheme?: Theme;
    themes?: Theme[];
    storageKey?: string;
}

interface ThemeManager {
    availableThemes: Theme[];
    theme: Ref<Theme>;
    setTheme: (theme: Theme) => void;
    setSystemTheme: () => void;
    toggleTheme: () => void;
    isDark: () => boolean;
}

const ThemeSymbol = Symbol('ThemeManager');

const availableThemes = [
    "light",
    "dark",
    "cupcake",
    "bumblebee",
    "emerald",
    "corporate",
    "synthwave",
    "retro",
    "cyberpunk",
    "valentine",
    "halloween",
    "garden",
    "forest",
    "aqua",
    "lofi",
    "pastel",
    "fantasy",
    "wireframe",
    "black",
    "luxury",
    "dracula",
    "cmyk",
    "autumn",
    "business",
    "acid",
    "lemonade",
    "night",
    "coffee",
    "winter",
    "dim",
    "nord",
    "sunset",
    "caramellatte",
    "abyss",
    "silk"
];
const darkThemes = [
    "dark",
    "cupcake",
    "bumblebee",
    "emerald",
    "corporate",
    "synthwave",
    "retro",
    "cyberpunk",
    "valentine",
    "halloween",
    "garden",
    "forest",
    "aqua",
    "lofi",
    "pastel",
    "fantasy",
    "wireframe",
    "black",
    "luxury",
    "dracula",
    "cmyk",
    "autumn",
    "business",
    "acid",
    "lemonade",
    "night",
    "coffee",
    "winter",
    "dim",
    "nord",
];


export function createThemePlugin(options: ThemePluginOptions = {}) {
    const themes = options.themes || availableThemes;
    const storageKey = options.storageKey ?? 'app-theme';
    const defaultTheme = options.defaultTheme ?? 'light';

    const storedTheme = localStorage.getItem(storageKey);
    const theme = ref<Theme>(storedTheme && themes.includes(storedTheme) ? storedTheme : defaultTheme);

    function setTheme(newTheme: Theme) {
        if (!themes.includes(newTheme)) {
            newTheme = defaultTheme;
        }

        theme.value = newTheme;
        localStorage.setItem(storageKey, newTheme);
        document.documentElement.setAttribute('data-theme', newTheme);
        document.documentElement.classList.remove("dark", "light");
        if (darkThemes.includes(newTheme)) {
            document.documentElement.classList.add("dark");
        } else {
            document.documentElement.classList.add("light");
        }
    }

    function setSystemTheme() {
        const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
        setTheme(systemTheme);
    }

    function toggleTheme() {
        const idx = themes.indexOf(theme.value);
        const nextTheme = themes[(idx + 1) % themes.length];
        setTheme(nextTheme);
    }

    function isDark(): boolean {
        return darkThemes.includes(theme.value);
    }

    // Set initial theme on load
    setTheme(theme.value);

    const manager: ThemeManager = {
        theme,
        availableThemes: themes,
        setTheme,
        setSystemTheme,
        toggleTheme,
        isDark,
    };

    return {
        install(app: App) {
            app.provide(ThemeSymbol, manager);
            app.config.globalProperties.$theme = manager;
        },
    };
}

export function useTheme(): ThemeManager {
    const manager = inject<ThemeManager>(ThemeSymbol);
    if (!manager) {
        throw new Error('ThemePlugin is not installed');
    }
    return manager;
}
