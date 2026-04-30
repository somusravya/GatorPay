import { Component, signal, inject, OnInit, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { WalletService } from '../../core/services/wallet.service';
import { Transaction } from '../../shared/models';

interface CategoryData {
    name: string;
    amount: number;
    color: string;
    icon: string;
    percentage: number;
}

interface MonthlyData {
    month: string;
    amount: number;
}

@Component({
    selector: 'app-spending-analytics',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './spending-analytics.component.html',
    styleUrl: './spending-analytics.component.scss'
})
export class SpendingAnalyticsComponent implements OnInit, AfterViewInit {
    @ViewChild('donutCanvas') donutCanvasRef!: ElementRef<HTMLCanvasElement>;
    @ViewChild('lineCanvas') lineCanvasRef!: ElementRef<HTMLCanvasElement>;
    @ViewChild('barCanvas') barCanvasRef!: ElementRef<HTMLCanvasElement>;

    private walletService = inject(WalletService);

    loading = signal(true);
    transactions = signal<Transaction[]>([]);
    categories = signal<CategoryData[]>([]);
    monthlyData = signal<MonthlyData[]>([]);
    totalSpent = signal(0);
    avgDaily = signal(0);
    topCategory = signal('');
    transactionCount = signal(0);
    selectedRange = signal<string>('3months');

    private categoryMap: Record<string, { color: string; icon: string; keywords: string[] }> = {
        'Food & Dining': { color: '#f59e0b', icon: '🍔', keywords: ['food', 'restaurant', 'dining', 'coffee', 'lunch', 'dinner', 'breakfast', 'meal'] },
        'Bills & Utilities': { color: '#3b82f6', icon: '📄', keywords: ['bill', 'utility', 'electric', 'water', 'gas', 'internet', 'phone', 'rent'] },
        'Shopping': { color: '#ec4899', icon: '🛍️', keywords: ['shopping', 'store', 'amazon', 'purchase', 'buy', 'order'] },
        'Transfers': { color: '#8b5cf6', icon: '💸', keywords: ['transfer', 'send', 'sent', 'payment'] },
        'Entertainment': { color: '#06b6d4', icon: '🎬', keywords: ['entertainment', 'movie', 'game', 'music', 'netflix', 'spotify', 'subscription'] },
        'Transport': { color: '#22c55e', icon: '🚗', keywords: ['uber', 'lyft', 'gas', 'fuel', 'transport', 'taxi', 'parking'] },
        'Health': { color: '#ef4444', icon: '🏥', keywords: ['health', 'medical', 'doctor', 'pharmacy', 'gym', 'fitness'] },
        'Other': { color: '#64748b', icon: '📦', keywords: [] }
    };

    ngOnInit(): void {
        this.loadTransactions();
    }

    ngAfterViewInit(): void {
        // Charts will render after data loads
    }

    setRange(range: string): void {
        this.selectedRange.set(range);
        this.processTransactions();
    }

    private loadTransactions(): void {
        this.walletService.getTransactions(1, 200).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.transactions.set(res.data.transactions || []);
                }
                this.processTransactions();
                this.loading.set(false);
            },
            error: () => {
                this.generateMockData();
                this.processTransactions();
                this.loading.set(false);
            }
        });
    }

    private generateMockData(): void {
        const types = ['withdraw', 'deposit'];
        const descriptions = [
            'Coffee Shop Purchase', 'Monthly Rent', 'Grocery Store', 'Netflix Subscription',
            'Uber Ride', 'Transfer to John', 'Restaurant Bill', 'Amazon Order',
            'Electric Bill', 'Gym Membership', 'Salary Deposit', 'Freelance Payment',
            'Gas Station', 'Movie Tickets', 'Phone Bill', 'Shopping Mall',
            'Doctor Visit', 'Music Subscription', 'Food Delivery', 'Water Bill'
        ];

        const mockTx: Transaction[] = [];
        for (let i = 0; i < 60; i++) {
            const date = new Date();
            date.setDate(date.getDate() - Math.floor(Math.random() * 180));
            const type = types[Math.floor(Math.random() * 2)];
            mockTx.push({
                id: `mock-${i}`,
                wallet_id: 'w1',
                from_user_id: 'u1',
                to_user_id: 'u2',
                type,
                amount: (Math.random() * 200 + 5).toFixed(2),
                description: descriptions[Math.floor(Math.random() * descriptions.length)],
                status: 'completed',
                created_at: date.toISOString()
            });
        }
        this.transactions.set(mockTx);
    }

    private processTransactions(): void {
        const range = this.selectedRange();
        const now = new Date();
        let cutoff = new Date();

        switch (range) {
            case '1month': cutoff.setMonth(now.getMonth() - 1); break;
            case '3months': cutoff.setMonth(now.getMonth() - 3); break;
            case '6months': cutoff.setMonth(now.getMonth() - 6); break;
            case 'all': cutoff = new Date(0); break;
        }

        const filtered = this.transactions().filter(tx => {
            const txDate = new Date(tx.created_at);
            return tx.type === 'withdraw' && txDate >= cutoff;
        });

        this.transactionCount.set(filtered.length);

        // Category breakdown
        const catTotals: Record<string, number> = {};
        for (const tx of filtered) {
            const cat = this.categorize(tx.description);
            catTotals[cat] = (catTotals[cat] || 0) + parseFloat(tx.amount);
        }

        const total = Object.values(catTotals).reduce((a, b) => a + b, 0);
        this.totalSpent.set(total);

        const days = Math.max(1, Math.ceil((now.getTime() - cutoff.getTime()) / (1000 * 60 * 60 * 24)));
        this.avgDaily.set(total / days);

        const cats: CategoryData[] = Object.entries(catTotals)
            .map(([name, amount]) => ({
                name,
                amount,
                color: this.categoryMap[name]?.color || '#64748b',
                icon: this.categoryMap[name]?.icon || '📦',
                percentage: total > 0 ? (amount / total) * 100 : 0
            }))
            .sort((a, b) => b.amount - a.amount);

        this.categories.set(cats);
        this.topCategory.set(cats[0]?.name || 'N/A');

        // Monthly data
        const monthTotals: Record<string, number> = {};
        for (const tx of filtered) {
            const d = new Date(tx.created_at);
            const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
            monthTotals[key] = (monthTotals[key] || 0) + parseFloat(tx.amount);
        }

        const monthly = Object.entries(monthTotals)
            .sort(([a], [b]) => a.localeCompare(b))
            .map(([month, amount]) => {
                const [y, m] = month.split('-');
                const date = new Date(parseInt(y), parseInt(m) - 1);
                return {
                    month: date.toLocaleDateString('en-US', { month: 'short', year: '2-digit' }),
                    amount
                };
            });
        this.monthlyData.set(monthly);

        // Render charts after a tick
        setTimeout(() => this.renderCharts(), 50);
    }

    private categorize(description: string): string {
        const lower = description.toLowerCase();
        for (const [cat, data] of Object.entries(this.categoryMap)) {
            if (cat === 'Other') continue;
            if (data.keywords.some(kw => lower.includes(kw))) return cat;
        }
        return 'Other';
    }

    private renderCharts(): void {
        this.renderDonut();
        this.renderLine();
        this.renderBar();
    }

    private renderDonut(): void {
        const canvas = this.donutCanvasRef?.nativeElement;
        if (!canvas) return;
        const ctx = canvas.getContext('2d')!;
        const dpr = window.devicePixelRatio || 1;
        const size = 220;
        canvas.width = size * dpr;
        canvas.height = size * dpr;
        canvas.style.width = size + 'px';
        canvas.style.height = size + 'px';
        ctx.scale(dpr, dpr);

        const cx = size / 2, cy = size / 2, radius = 85, lineWidth = 28;
        const cats = this.categories();
        const total = cats.reduce((a, c) => a + c.amount, 0);

        ctx.clearRect(0, 0, size, size);

        // Background ring
        ctx.beginPath();
        ctx.arc(cx, cy, radius, 0, Math.PI * 2);
        ctx.strokeStyle = 'rgba(100,116,139,0.15)';
        ctx.lineWidth = lineWidth;
        ctx.stroke();

        // Segments
        let startAngle = -Math.PI / 2;
        for (const cat of cats) {
            const sweep = total > 0 ? (cat.amount / total) * Math.PI * 2 : 0;
            ctx.beginPath();
            ctx.arc(cx, cy, radius, startAngle, startAngle + sweep);
            ctx.strokeStyle = cat.color;
            ctx.lineWidth = lineWidth;
            ctx.lineCap = 'round';
            ctx.stroke();
            startAngle += sweep + 0.03;
        }

        // Center text
        ctx.fillStyle = getComputedStyle(document.documentElement).getPropertyValue('--text-primary').trim() || '#f1f5f9';
        ctx.font = 'bold 20px Inter';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText('$' + total.toFixed(0), cx, cy - 8);
        ctx.font = '12px Inter';
        ctx.fillStyle = getComputedStyle(document.documentElement).getPropertyValue('--text-muted').trim() || '#64748b';
        ctx.fillText('Total Spent', cx, cy + 12);
    }

    private renderLine(): void {
        const canvas = this.lineCanvasRef?.nativeElement;
        if (!canvas) return;
        const ctx = canvas.getContext('2d')!;
        const dpr = window.devicePixelRatio || 1;
        const w = 500, h = 200;
        canvas.width = w * dpr;
        canvas.height = h * dpr;
        canvas.style.width = '100%';
        canvas.style.height = h + 'px';
        ctx.scale(dpr, dpr);

        const data = this.monthlyData();
        if (data.length === 0) return;

        const padding = { top: 20, right: 20, bottom: 35, left: 55 };
        const chartW = w - padding.left - padding.right;
        const chartH = h - padding.top - padding.bottom;
        const maxVal = Math.max(...data.map(d => d.amount), 1);

        ctx.clearRect(0, 0, w, h);

        // Grid lines
        const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-muted').trim() || '#64748b';
        for (let i = 0; i <= 4; i++) {
            const y = padding.top + (chartH / 4) * i;
            ctx.strokeStyle = 'rgba(100,116,139,0.1)';
            ctx.lineWidth = 1;
            ctx.beginPath();
            ctx.moveTo(padding.left, y);
            ctx.lineTo(w - padding.right, y);
            ctx.stroke();

            ctx.fillStyle = textColor;
            ctx.font = '11px Inter';
            ctx.textAlign = 'right';
            ctx.fillText('$' + Math.round(maxVal - (maxVal / 4) * i), padding.left - 8, y + 4);
        }

        // Line
        const gradient = ctx.createLinearGradient(0, padding.top, 0, h - padding.bottom);
        gradient.addColorStop(0, 'rgba(59,130,246,0.3)');
        gradient.addColorStop(1, 'rgba(59,130,246,0)');

        const points: { x: number; y: number }[] = data.map((d, i) => ({
            x: padding.left + (chartW / Math.max(data.length - 1, 1)) * i,
            y: padding.top + chartH - (d.amount / maxVal) * chartH
        }));

        // Fill area
        ctx.beginPath();
        ctx.moveTo(points[0].x, h - padding.bottom);
        points.forEach(p => ctx.lineTo(p.x, p.y));
        ctx.lineTo(points[points.length - 1].x, h - padding.bottom);
        ctx.closePath();
        ctx.fillStyle = gradient;
        ctx.fill();

        // Line stroke
        ctx.beginPath();
        points.forEach((p, i) => i === 0 ? ctx.moveTo(p.x, p.y) : ctx.lineTo(p.x, p.y));
        ctx.strokeStyle = '#3b82f6';
        ctx.lineWidth = 2.5;
        ctx.lineJoin = 'round';
        ctx.stroke();

        // Dots
        points.forEach(p => {
            ctx.beginPath();
            ctx.arc(p.x, p.y, 4, 0, Math.PI * 2);
            ctx.fillStyle = '#3b82f6';
            ctx.fill();
            ctx.strokeStyle = getComputedStyle(document.documentElement).getPropertyValue('--bg-card').trim() || '#1e293b';
            ctx.lineWidth = 2;
            ctx.stroke();
        });

        // X labels
        ctx.fillStyle = textColor;
        ctx.font = '11px Inter';
        ctx.textAlign = 'center';
        data.forEach((d, i) => {
            const x = padding.left + (chartW / Math.max(data.length - 1, 1)) * i;
            ctx.fillText(d.month, x, h - 10);
        });
    }

    private renderBar(): void {
        const canvas = this.barCanvasRef?.nativeElement;
        if (!canvas) return;
        const ctx = canvas.getContext('2d')!;
        const dpr = window.devicePixelRatio || 1;
        const w = 500, h = 200;
        canvas.width = w * dpr;
        canvas.height = h * dpr;
        canvas.style.width = '100%';
        canvas.style.height = h + 'px';
        ctx.scale(dpr, dpr);

        const cats = this.categories().slice(0, 6);
        if (cats.length === 0) return;

        const padding = { top: 20, right: 20, bottom: 40, left: 55 };
        const chartW = w - padding.left - padding.right;
        const chartH = h - padding.top - padding.bottom;
        const maxVal = Math.max(...cats.map(c => c.amount), 1);
        const barW = Math.min(40, chartW / cats.length - 12);

        ctx.clearRect(0, 0, w, h);

        const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-muted').trim() || '#64748b';

        // Grid
        for (let i = 0; i <= 4; i++) {
            const y = padding.top + (chartH / 4) * i;
            ctx.strokeStyle = 'rgba(100,116,139,0.1)';
            ctx.lineWidth = 1;
            ctx.beginPath();
            ctx.moveTo(padding.left, y);
            ctx.lineTo(w - padding.right, y);
            ctx.stroke();
        }

        // Bars
        cats.forEach((cat, i) => {
            const x = padding.left + (chartW / cats.length) * i + (chartW / cats.length - barW) / 2;
            const barH = (cat.amount / maxVal) * chartH;
            const y = padding.top + chartH - barH;

            // Bar with rounded top
            const radius = 5;
            ctx.beginPath();
            ctx.moveTo(x, y + radius);
            ctx.quadraticCurveTo(x, y, x + radius, y);
            ctx.lineTo(x + barW - radius, y);
            ctx.quadraticCurveTo(x + barW, y, x + barW, y + radius);
            ctx.lineTo(x + barW, padding.top + chartH);
            ctx.lineTo(x, padding.top + chartH);
            ctx.closePath();
            ctx.fillStyle = cat.color;
            ctx.fill();

            // Label
            ctx.fillStyle = textColor;
            ctx.font = '10px Inter';
            ctx.textAlign = 'center';
            ctx.fillText(cat.icon, x + barW / 2, h - 8);
        });
    }
}
