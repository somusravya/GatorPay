import { Injectable, signal, effect } from '@angular/core';

export type ThemeMode = 'dark' | 'light' | 'system';

@Injectable({ providedIn: 'root' })
export class ThemeService {
    private readonly storageKey = 'gatorpay_theme';
    private mediaQuery = globalThis.matchMedia?.('(prefers-color-scheme: dark)');

    theme = signal<ThemeMode>(this.loadTheme());
    isDark = signal<boolean>(true);

    constructor() {
        // Apply theme whenever the signal changes
        effect(() => {
            this.applyTheme(this.theme());
        });

        // Listen for system theme changes
        this.mediaQuery?.addEventListener('change', () => {
            if (this.theme() === 'system') {
                this.applyTheme('system');
            }
        });
    }

    toggle(): void {
        const modes: ThemeMode[] = ['dark', 'light', 'system'];
        const currentIndex = modes.indexOf(this.theme());
        const next = modes[(currentIndex + 1) % modes.length];
        this.theme.set(next);
        localStorage.setItem(this.storageKey, next);
    }

    setTheme(mode: ThemeMode): void {
        this.theme.set(mode);
        localStorage.setItem(this.storageKey, mode);
    }

    private loadTheme(): ThemeMode {
        const stored = localStorage.getItem(this.storageKey) as ThemeMode | null;
        return stored || 'dark';
    }

    private applyTheme(mode: ThemeMode): void {
        let resolved: 'dark' | 'light';

        if (mode === 'system') {
            resolved = this.mediaQuery?.matches ? 'dark' : 'light';
        } else {
            resolved = mode;
        }

        this.isDark.set(resolved === 'dark');
        document.documentElement.setAttribute('data-theme', resolved);
    }

    getIcon(): string {
        switch (this.theme()) {
            case 'dark': return '🌙';
            case 'light': return '☀️';
            case 'system': return '💻';
        }
    }

    getLabel(): string {
        switch (this.theme()) {
            case 'dark': return 'Dark';
            case 'light': return 'Light';
            case 'system': return 'System';
        }
    }
}
