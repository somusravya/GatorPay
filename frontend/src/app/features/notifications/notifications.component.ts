import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { NotificationService } from '../../core/services/notification.service';

@Component({
    selector: 'app-notifications',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './notifications.component.html',
    styleUrl: './notifications.component.scss'
})
export class NotificationsComponent implements OnInit {
    private notifService = inject(NotificationService);

    notifications = signal<any[]>([]);
    loading = signal(true);
    activeFilter = 'all';
    showPreferences = false;

    filters = ['all', 'payment', 'alert', 'promo', 'system'];

    preferences: Record<string, boolean> = {
        payment_alerts: true,
        trade_alerts: true,
        loan_reminders: true,
        promo_offers: true,
        security_alerts: true,
        price_alerts: false
    };

    preferenceLabels: { [key: string]: { label: string; icon: string; desc: string } } = {
        payment_alerts: { label: 'Payment Alerts', icon: '💳', desc: 'Sent/received payment notifications' },
        trade_alerts: { label: 'Trade Alerts', icon: '📈', desc: 'Order executions & portfolio updates' },
        loan_reminders: { label: 'Loan Reminders', icon: '🏦', desc: 'Payment due dates & status changes' },
        promo_offers: { label: 'Promo Offers', icon: '🎁', desc: 'Special deals & reward opportunities' },
        security_alerts: { label: 'Security Alerts', icon: '🛡️', desc: 'Login attempts & suspicious activity' },
        price_alerts: { label: 'Price Alerts', icon: '📊', desc: 'Stock price target notifications' }
    };

    ngOnInit(): void {
        this.loadNotifications('all');
        this.loadPreferences();
    }

    loadNotifications(filter: string): void {
        this.activeFilter = filter;
        this.loading.set(true);
        const typeFilter = filter === 'all' ? undefined : filter;
        this.notifService.getNotifications(typeFilter).subscribe({
            next: (res: any) => {
                if (res.success) {
                    const payload = res.data;
                    this.notifications.set(Array.isArray(payload) ? payload : (payload?.notifications || []));
                }
                this.loading.set(false);
            },
            error: () => this.loading.set(false)
        });
    }

    loadPreferences(): void {
        this.notifService.getPreferences().subscribe({
            next: (res: any) => { if (res.success && res.data) this.preferences = { ...this.preferences, ...res.data }; }
        });
    }

    savePreferences(): void {
        this.notifService.updatePreferences(this.preferences).subscribe({
            next: () => { this.showPreferences = false; }
        });
    }

    preferenceKeys(): string[] {
        return Object.keys(this.preferenceLabels);
    }

    unreadCount(): number {
        return this.notifications().filter(n => !n.is_read).length;
    }

    markRead(notif: any): void {
        if (notif.is_read) return;
        this.notifService.markAsRead(notif.id).subscribe({
            next: () => { notif.is_read = true; }
        });
    }

    markAllRead(): void {
        this.notifService.markAllRead().subscribe({
            next: () => {
                this.notifications.update(list => list.map(n => ({ ...n, is_read: true })));
            }
        });
    }

    getTypeIcon(type: string): string {
        const icons: { [key: string]: string } = {
            payment: '💳', alert: '⚠️', promo: '🎁', system: '⚙️', trade: '📈', loan: '🏦'
        };
        return icons[type] || '🔔';
    }

    getTimeAgo(dateStr: string): string {
        const diff = Date.now() - new Date(dateStr).getTime();
        const mins = Math.floor(diff / 60000);
        if (mins < 60) return `${mins}m ago`;
        const hours = Math.floor(mins / 60);
        if (hours < 24) return `${hours}h ago`;
        return `${Math.floor(hours / 24)}d ago`;
    }
}
