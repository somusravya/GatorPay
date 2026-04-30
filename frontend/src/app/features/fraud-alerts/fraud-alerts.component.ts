import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { FraudService } from '../../core/services/fraud.service';

@Component({
    selector: 'app-fraud-alerts',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './fraud-alerts.component.html',
    styleUrl: './fraud-alerts.component.scss'
})
export class FraudAlertsComponent implements OnInit {
    private fraudService = inject(FraudService);

    alerts = signal<any[]>([]);
    loading = signal(true);
    activeFilter = 'all';
    showDetail = signal<any>(null);
    actionLoading = signal('');

    filters = [
        { id: 'all', label: 'All Alerts', icon: '📋' },
        { id: 'pending', label: 'Pending', icon: '⏳' },
        { id: 'confirmed', label: 'Confirmed', icon: '✅' },
        { id: 'dismissed', label: 'Dismissed', icon: '❌' }
    ];

    ngOnInit(): void {
        this.loadAlerts();
    }

    loadAlerts(): void {
        this.loading.set(true);
        const status = this.activeFilter === 'all' ? undefined : this.activeFilter;
        this.fraudService.getAlerts(status).subscribe({
            next: (res: any) => {
                if (res.success) this.alerts.set(res.data || []);
                this.loading.set(false);
            },
            error: () => this.loading.set(false)
        });
    }

    filterAlerts(filter: string): void {
        this.activeFilter = filter;
        this.loadAlerts();
    }

    reviewAlert(alertId: string, action: string): void {
        this.actionLoading.set(alertId);
        this.fraudService.reviewAlert({ alert_id: alertId, action }).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.alerts.update(list =>
                        list.map(a => a.id === alertId ? { ...a, status: action } : a)
                    );
                }
                this.actionLoading.set('');
                this.showDetail.set(null);
            },
            error: () => this.actionLoading.set('')
        });
    }

    get pendingCount(): number {
        return this.alerts().filter(a => a.status === 'pending').length;
    }

    get highRiskCount(): number {
        return this.alerts().filter(a => a.risk_score >= 70).length;
    }

    get resolvedCount(): number {
        return this.alerts().filter(a => a.status !== 'pending').length;
    }

    getRiskLevel(score: number): string {
        if (score >= 70) return 'critical';
        if (score >= 40) return 'medium';
        return 'low';
    }

    getRiskLabel(score: number): string {
        if (score >= 70) return 'High Risk';
        if (score >= 40) return 'Medium Risk';
        return 'Low Risk';
    }

    getTypeIcon(type: string): string {
        const icons: Record<string, string> = {
            velocity: '⚡', amount_spike: '📈', geo_anomaly: '🌍', suspicious_merchant: '🏪'
        };
        return icons[type] || '⚠️';
    }

    getTypeLabel(type: string): string {
        const labels: Record<string, string> = {
            velocity: 'Rapid Transactions', amount_spike: 'Unusual Amount',
            geo_anomaly: 'Unusual Time', suspicious_merchant: 'Suspicious Activity'
        };
        return labels[type] || 'Unknown';
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
