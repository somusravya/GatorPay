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

    billTotal = 96;
    splitPeople = 4;
    tipPercent = 18;

    loanAmount = 6500;
    loanApr = 9.5;
    loanMonths = 24;

    emergencyMonthlySpend = 2200;
    emergencySaved = 1800;

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

    getSplitTotal(): number {
        return this.billTotal + (this.billTotal * this.tipPercent / 100);
    }

    getSplitAmount(): number {
        if (this.splitPeople <= 0) return 0;
        return this.getSplitTotal() / this.splitPeople;
    }

    getMonthlyLoanPayment(): number {
        if (this.loanMonths <= 0) return 0;
        const monthlyRate = this.loanApr / 100 / 12;
        if (monthlyRate === 0) return this.loanAmount / this.loanMonths;
        const factor = Math.pow(1 + monthlyRate, this.loanMonths);
        return this.loanAmount * monthlyRate * factor / (factor - 1);
    }

    getTotalLoanInterest(): number {
        return Math.max((this.getMonthlyLoanPayment() * this.loanMonths) - this.loanAmount, 0);
    }

    getEmergencyTarget(): number {
        return this.emergencyMonthlySpend * 3;
    }

    getEmergencyProgress(): number {
        const target = this.getEmergencyTarget();
        if (target <= 0) return 0;
        return Math.min((this.emergencySaved / target) * 100, 100);
    }

    getEmergencyGap(): number {
        return Math.max(this.getEmergencyTarget() - this.emergencySaved, 0);
    }
}
