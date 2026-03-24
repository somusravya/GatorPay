import { Component, computed, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';

@Component({
    selector: 'app-profile',
    standalone: true,
    imports: [CommonModule, RouterLink],
    templateUrl: './profile.component.html',
    styleUrl: './profile.component.scss'
})
export class ProfileComponent {
    private authService = inject(AuthService);

    user = this.authService.currentUser;
    wallet = this.authService.currentWallet;

    userInitials = computed(() => {
        const u = this.user();
        if (!u) return '?';
        return (u.first_name[0] + u.last_name[0]).toUpperCase();
    });

    fullName = computed(() => {
        const u = this.user();
        if (!u) return '';
        return `${u.first_name} ${u.last_name}`;
    });

    memberSince = computed(() => {
        const u = this.user();
        if (!u) return '';
        return new Date(u.created_at).toLocaleDateString('en-US', {
            year: 'numeric', month: 'long', day: 'numeric'
        });
    });

    getFormattedBalance(): string {
        const balance = this.wallet()?.balance || '0';
        return parseFloat(balance).toLocaleString('en-US', { style: 'currency', currency: 'USD' });
    }

    logout(): void {
        this.authService.logout();
    }
}
