import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AdminService } from '../../core/services/admin.service';

@Component({
    selector: 'app-admin',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './admin.component.html',
    styleUrl: './admin.component.scss'
})
export class AdminComponent implements OnInit {
    private adminService = inject(AdminService);
    metrics = signal<any>(null);
    users = signal<any[]>([]);
    fraudAlerts = signal<any[]>([]);
    loading = signal(true);
    activeTab = 'overview';
    searchQuery = '';
    userPage = 1;
    userTotal = signal(0);

    tabs = [
        { id: 'overview', label: 'Overview', icon: '📊' },
        { id: 'users', label: 'Users', icon: '👥' },
        { id: 'fraud', label: 'Fraud Review', icon: '🛡️' },
    ];

    ngOnInit(): void {
        this.loadMetrics();
        this.loadUsers();
        this.loadFraudQueue();
    }

    loadMetrics(): void {
        this.adminService.getMetrics().subscribe({
            next: (res: any) => { if (res.success) this.metrics.set(res.data); this.loading.set(false); },
            error: () => this.loading.set(false)
        });
    }

    loadUsers(): void {
        this.adminService.getUsers(this.userPage, 20, this.searchQuery).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.users.set(res.data?.users || []);
                    this.userTotal.set(res.data?.total || 0);
                }
            }
        });
    }

    loadFraudQueue(): void {
        this.adminService.getFraudReview().subscribe({
            next: (res: any) => { if (res.success) this.fraudAlerts.set(res.data || []); }
        });
    }

    searchUsers(): void {
        this.userPage = 1;
        this.loadUsers();
    }

    setTab(tab: string): void { this.activeTab = tab; }
}
