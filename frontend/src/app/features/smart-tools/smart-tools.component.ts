import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
    selector: 'app-smart-tools',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './smart-tools.component.html',
    styleUrl: './smart-tools.component.scss'
})
export class SmartToolsComponent {
    monthlyIncome = 4200;
    monthlySpending = 2850;
    fixedBills = 1250;
    savingsTarget = 500;

    selectedCategory = 'Food';
    categorySpent = 185;
    categoryLimit = 450;

    goalName = 'Spring break fund';
    goalAmount = 2500;
    currentSavings = 820;
    monthlyContribution = 260;

    getRemainingAfterSpending(): number {
        return Math.max(this.monthlyIncome - this.monthlySpending, 0);
    }

    getFlexibleBudget(): number {
        return Math.max(this.monthlyIncome - this.fixedBills - this.savingsTarget, 0);
    }

    getSavingsRate(): number {
        if (this.monthlyIncome <= 0) return 0;
        return Math.min((this.savingsTarget / this.monthlyIncome) * 100, 100);
    }

    getBudgetHealth(): string {
        const remaining = this.getRemainingAfterSpending();
        if (remaining >= this.savingsTarget) return 'On track';
        if (remaining > 0) return 'Needs review';
        return 'Over budget';
    }

    getCategoryUsage(): number {
        if (this.categoryLimit <= 0) return 0;
        return Math.min((this.categorySpent / this.categoryLimit) * 100, 100);
    }

    getCategoryStatus(): string {
        const usage = this.getCategoryUsage();
        if (usage >= 90) return 'Limit warning';
        if (usage >= 70) return 'Watch closely';
        return 'Healthy';
    }

    getGoalProgress(): number {
        if (this.goalAmount <= 0) return 0;
        return Math.min((this.currentSavings / this.goalAmount) * 100, 100);
    }

    getMonthsToGoal(): number {
        const remaining = Math.max(this.goalAmount - this.currentSavings, 0);
        if (remaining === 0) return 0;
        if (this.monthlyContribution <= 0) return 0;
        return Math.ceil(remaining / this.monthlyContribution);
    }
}
