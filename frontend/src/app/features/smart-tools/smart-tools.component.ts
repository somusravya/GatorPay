import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

interface SubscriptionItem {
    name: string;
    cost: number;
}

interface FraudSignal {
    label: string;
    checked: boolean;
    weight: number;
}

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

    subscriptions: SubscriptionItem[] = [
        { name: 'Streaming bundle', cost: 22 },
        { name: 'Cloud storage', cost: 9 },
        { name: 'Fitness app', cost: 14 }
    ];

    cashbackPurchaseAmount = 350;
    cashbackRate = 2;

    fraudSignals: FraudSignal[] = [
        { label: 'New recipient or merchant', checked: false, weight: 25 },
        { label: 'Amount is higher than usual', checked: false, weight: 20 },
        { label: 'Payment requested urgently', checked: false, weight: 25 },
        { label: 'Link or QR came from an unknown sender', checked: false, weight: 30 }
    ];

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

    getMonthlySubscriptionTotal(): number {
        return this.subscriptions.reduce((total, item) => total + Number(item.cost || 0), 0);
    }

    getAnnualSubscriptionTotal(): number {
        return this.getMonthlySubscriptionTotal() * 12;
    }

    addSubscription(): void {
        this.subscriptions = [
            ...this.subscriptions,
            { name: `Subscription ${this.subscriptions.length + 1}`, cost: 10 }
        ];
    }

    removeSubscription(index: number): void {
        this.subscriptions = this.subscriptions.filter((_, itemIndex) => itemIndex !== index);
    }

    getCashbackEstimate(): number {
        return this.cashbackPurchaseAmount * (this.cashbackRate / 100);
    }

    getNetAfterCashback(): number {
        return Math.max(this.cashbackPurchaseAmount - this.getCashbackEstimate(), 0);
    }

    toggleFraudSignal(index: number): void {
        this.fraudSignals = this.fraudSignals.map((item, itemIndex) => (
            itemIndex === index ? { ...item, checked: !item.checked } : item
        ));
    }

    getFraudRiskScore(): number {
        return this.fraudSignals.reduce((score, item) => score + (item.checked ? item.weight : 0), 0);
    }

    getFraudRiskLabel(): string {
        const score = this.getFraudRiskScore();
        if (score >= 70) return 'High risk';
        if (score >= 35) return 'Review carefully';
        return 'Looks normal';
    }

    getSmartActions(): string[] {
        const savingsMove = Math.min(this.getRemainingAfterSpending(), this.savingsTarget);
        return [
            `Move $${savingsMove.toFixed(0)} toward savings this month.`,
            `Keep ${this.selectedCategory} below ${this.categoryLimit.toFixed(0)} for this cycle.`,
            `Review $${this.getMonthlySubscriptionTotal().toFixed(0)} in monthly subscriptions before the next renewal.`,
            `Emergency fund gap is $${this.getEmergencyGap().toFixed(0)} toward a 3-month cushion.`
        ];
    }
}
