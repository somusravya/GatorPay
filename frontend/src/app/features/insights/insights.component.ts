import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { InsightService } from '../../core/services/insight.service';

@Component({
    selector: 'app-insights',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './insights.component.html',
    styleUrl: './insights.component.scss'
})
export class InsightsComponent implements OnInit {
    private insightService = inject(InsightService);

    loading = signal(true);
    insights = signal<any>(null);
    error = signal('');

    ngOnInit(): void {
        this.loadInsights();
    }

    loadInsights(): void {
        this.insightService.getSummary().subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.insights.set(res.data);
                }
                this.loading.set(false);
            },
            error: () => {
                this.error.set('Failed to load insights');
                this.loading.set(false);
            }
        });
    }

    getGradeColor(grade: string): string {
        const colors: Record<string, string> = {
            'Excellent': '#22c55e',
            'Good': '#3b82f6',
            'Fair': '#f59e0b',
            'Needs Improvement': '#f97316',
            'Critical': '#ef4444'
        };
        return colors[grade] || '#94a3b8';
    }

    getScoreArc(score: number): string {
        const radius = 80;
        const circumference = 2 * Math.PI * radius;
        const offset = circumference - (score / 100) * circumference;
        return `${circumference} ${offset}`;
    }

    Math = Math;
}
