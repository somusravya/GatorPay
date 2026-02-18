import { Component, computed, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
    selector: 'app-layout',
    standalone: true,
    imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive],
    templateUrl: './layout.component.html',
    styleUrl: './layout.component.scss'
})
export class LayoutComponent {
    private authService = inject(AuthService);

    sidebarCollapsed = signal(false);
    user = this.authService.currentUser;
    wallet = this.authService.currentWallet;

    greeting = computed(() => {
        const hour = new Date().getHours();
        if (hour < 12) return 'Good Morning';
        if (hour < 17) return 'Good Afternoon';
        return 'Good Evening';
    });

    userInitials = computed(() => {
        const u = this.user();
        if (!u) return '?';
        return (u.first_name[0] + u.last_name[0]).toUpperCase();
    });

    kycBadgeClass = computed(() => {
        const status = this.user()?.kyc_status || 'pending';
        return `kyc-badge kyc-${status}`;
    });

    toggleSidebar(): void {
        this.sidebarCollapsed.update(v => !v);
    }

    logout(): void {
        this.authService.logout();
    }
}
