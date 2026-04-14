import { Component, signal, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RewardService } from '../../core/services/reward.service';
import { RewardSummary, Reward, Offer } from '../../shared/models';

@Component({
    selector: 'app-rewards',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './rewards.component.html',
    styleUrl: './rewards.component.scss'
})
export class RewardsComponent implements OnInit {
    private rewardService = inject(RewardService);

    summary = signal<RewardSummary | null>(null);
    rewards = signal<Reward[]>([]);
    offers = signal<Offer[]>([]);
    loading = signal(true);
    historyPage = 1;
    historyTotal = signal(0);

    ngOnInit(): void {
        this.loadSummary();
        this.loadHistory();
        this.loadOffers();
    }

    loadSummary(): void {
        this.rewardService.getSummary().subscribe({
            next: (res) => {
                if (res.success) {
                    this.summary.set(res.data);
                }
                this.loading.set(false);
            },
            error: () => this.loading.set(false)
        });
    }

    loadHistory(): void {
        this.rewardService.getHistory(this.historyPage, 10).subscribe({
            next: (res) => {
                if (res.success) {
                    this.rewards.set(res.data.rewards || []);
                    this.historyTotal.set(res.data.total);
                }
            }
        });
    }

    loadOffers(): void {
        this.rewardService.getOffers().subscribe({
            next: (res) => {
                if (res.success) {
                    this.offers.set(res.data || []);
                }
            }
        });
    }

    formatCurrency(val: string): string {
        return parseFloat(val || '0').toLocaleString('en-US', { style: 'currency', currency: 'USD' });
    }

    nextPage(): void {
        this.historyPage++;
        this.loadHistory();
    }

    prevPage(): void {
        if (this.historyPage > 1) {
            this.historyPage--;
            this.loadHistory();
        }
    }
}
