import { Component, OnInit, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SavingsGoalsService, SavingsGoal } from '../../core/services/savings-goals.service';

@Component({
    selector: 'app-savings-goals',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './savings-goals.component.html',
    styleUrl: './savings-goals.component.scss'
})
export class SavingsGoalsComponent implements OnInit {
    private goalsService = inject(SavingsGoalsService);

    goals = signal<SavingsGoal[]>([]);
    showCreateModal = signal(false);
    showAddMoneyModal = signal(false);
    selectedGoalId = signal<string | null>(null);
    addAmount = 0;
    celebratingGoalId = signal<string | null>(null);

    newGoal = {
        name: '',
        targetAmount: 0,
        icon: '🎯',
        color: '#3b82f6',
        deadline: ''
    };

    iconOptions = ['🎯', '✈️', '🏠', '🚗', '💻', '📱', '🎓', '💍', '🏖️', '🎮', '🎸', '💪'];
    colorOptions = ['#3b82f6', '#8b5cf6', '#ec4899', '#22c55e', '#f59e0b', '#06b6d4', '#ef4444', '#6366f1'];

    ngOnInit(): void {
        this.loadGoals();
    }

    loadGoals(): void {
        this.goalsService.getGoals().subscribe(goals => this.goals.set(goals));
    }

    getProgress(goal: SavingsGoal): number {
        if (goal.targetAmount <= 0) return 0;
        return Math.min((goal.savedAmount / goal.targetAmount) * 100, 100);
    }

    getStrokeDasharray(goal: SavingsGoal): string {
        const circumference = 2 * Math.PI * 54;
        const progress = this.getProgress(goal) / 100;
        return `${circumference * progress} ${circumference * (1 - progress)}`;
    }

    getCircumference(): number {
        return 2 * Math.PI * 54;
    }

    getDaysLeft(goal: SavingsGoal): string {
        if (!goal.deadline) return 'No deadline';
        const diff = new Date(goal.deadline).getTime() - Date.now();
        if (diff <= 0) return 'Overdue';
        const days = Math.ceil(diff / (1000 * 60 * 60 * 24));
        return days === 1 ? '1 day left' : `${days} days left`;
    }

    openCreateModal(): void {
        this.newGoal = { name: '', targetAmount: 0, icon: '🎯', color: '#3b82f6', deadline: '' };
        this.showCreateModal.set(true);
    }

    closeCreateModal(): void {
        this.showCreateModal.set(false);
    }

    createGoal(): void {
        if (!this.newGoal.name.trim() || this.newGoal.targetAmount <= 0) return;
        this.goalsService.createGoal(this.newGoal).subscribe(() => {
            this.loadGoals();
            this.closeCreateModal();
        });
    }

    openAddMoney(goalId: string): void {
        this.selectedGoalId.set(goalId);
        this.addAmount = 0;
        this.showAddMoneyModal.set(true);
    }

    closeAddMoney(): void {
        this.showAddMoneyModal.set(false);
        this.selectedGoalId.set(null);
    }

    confirmAddMoney(): void {
        const goalId = this.selectedGoalId();
        if (!goalId || this.addAmount <= 0) return;

        const goalBefore = this.goals().find(g => g.id === goalId);
        const prevProgress = goalBefore ? this.getProgress(goalBefore) : 0;

        this.goalsService.addMoney(goalId, this.addAmount).subscribe(() => {
            this.loadGoals();
            this.closeAddMoney();

            // Check for milestone celebration
            const goalAfter = this.goals().find(g => g.id === goalId);
            if (goalAfter) {
                const newProgress = this.getProgress(goalAfter);
                const milestones = [25, 50, 75, 100];
                for (const m of milestones) {
                    if (prevProgress < m && newProgress >= m) {
                        this.celebrate(goalId);
                        break;
                    }
                }
            }
        });
    }

    deleteGoal(goalId: string): void {
        this.goalsService.deleteGoal(goalId).subscribe(() => this.loadGoals());
    }

    private celebrate(goalId: string): void {
        this.celebratingGoalId.set(goalId);
        setTimeout(() => this.celebratingGoalId.set(null), 2500);
    }

    selectIcon(icon: string): void {
        this.newGoal.icon = icon;
    }

    selectColor(color: string): void {
        this.newGoal.color = color;
    }
}
