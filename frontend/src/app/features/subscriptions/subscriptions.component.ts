import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { SubscriptionService } from '../../core/services/subscription.service';

@Component({
    selector: 'app-subscriptions',
    standalone: true,
    imports: [CommonModule, ReactiveFormsModule, FormsModule],
    templateUrl: './subscriptions.component.html',
    styleUrl: './subscriptions.component.scss'
})
export class SubscriptionsComponent implements OnInit {
    private subService = inject(SubscriptionService);
    private fb = inject(FormBuilder);

    subscriptions = signal<any[]>([]);
    loading = signal(true);
    showTrackForm = false;
    trackForm: FormGroup;

    categories = ['streaming', 'music', 'cloud', 'fitness', 'software', 'gaming', 'news', 'other'];

    constructor() {
        this.trackForm = this.fb.group({
            name: ['', Validators.required],
            category: ['streaming', Validators.required],
            amount: [9.99, [Validators.required, Validators.min(0.01)]],
            frequency: ['monthly', Validators.required],
            provider: [''],
            icon: ['🔄'],
            color: ['#8b5cf6']
        });
    }

    ngOnInit(): void {
        this.subService.getSubscriptions().subscribe({
            next: (res: any) => {
                if (res.success) this.subscriptions.set(res.data || []);
                this.loading.set(false);
            },
            error: () => this.loading.set(false)
        });
    }

    get totalMonthly(): number {
        return this.subscriptions().reduce((sum, s) => sum + parseFloat(s.amount || 0), 0);
    }

    get activeCount(): number {
        return this.subscriptions().filter(s => s.status === 'active').length;
    }

    get upcomingRenewals(): any[] {
        const now = new Date();
        const twoWeeks = new Date(now.getTime() + 14 * 24 * 60 * 60 * 1000);
        return this.subscriptions().filter(s => {
            if (!s.next_renewal) return false;
            const renewal = new Date(s.next_renewal);
            return renewal >= now && renewal <= twoWeeks;
        });
    }

    get optimizationSuggestions(): any[] {
        return [
            { icon: '💡', title: 'Bundle Savings', desc: 'Consider Disney+/Hulu/ESPN+ bundle to save $8/mo', savings: 8 },
            { icon: '🔄', title: 'Annual Plan Switch', desc: 'Switch Spotify to annual billing and save $20/yr', savings: 1.67 },
            { icon: '⚠️', title: 'Unused Service', desc: 'Adobe CC hasn\'t been used in 30 days — consider pausing', savings: 54.99 }
        ];
    }

    toggleAutoPay(sub: any): void {
        this.subService.setAutoPay({ subscription_id: sub.id, auto_pay: !sub.auto_pay }).subscribe({
            next: () => { sub.auto_pay = !sub.auto_pay; }
        });
    }

    trackSubscription(): void {
        if (this.trackForm.invalid) return;
        this.subService.trackSubscription(this.trackForm.value).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.subscriptions.update(s => [res.data, ...s]);
                    this.showTrackForm = false;
                    this.trackForm.reset({ category: 'streaming', amount: 9.99, frequency: 'monthly', icon: '🔄', color: '#8b5cf6' });
                }
            }
        });
    }

    getDaysUntilRenewal(dateStr: string): number {
        const renewal = new Date(dateStr);
        const now = new Date();
        return Math.ceil((renewal.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
    }
}
