import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { BudgetService } from '../../core/services/budget.service';

@Component({
    selector: 'app-budget-planner',
    standalone: true,
    imports: [CommonModule, ReactiveFormsModule, FormsModule],
    templateUrl: './budget-planner.component.html',
    styleUrl: './budget-planner.component.scss'
})
export class BudgetPlannerComponent implements OnInit {
    private budgetService = inject(BudgetService);
    private fb = inject(FormBuilder);

    goals = signal<any[]>([]);
    loading = signal(true);
    showGoalForm = signal(false);
    goalForm: FormGroup;

    categories = [
        { value: 'emergency', label: 'Emergency Fund', icon: '🛡️' },
        { value: 'vacation', label: 'Vacation', icon: '✈️' },
        { value: 'education', label: 'Education', icon: '📚' },
        { value: 'car', label: 'Car', icon: '🚗' },
        { value: 'home', label: 'Home', icon: '🏠' },
        { value: 'retirement', label: 'Retirement', icon: '🏖️' },
        { value: 'wedding', label: 'Wedding', icon: '💍' },
        { value: 'other', label: 'Other', icon: '🎯' }
    ];

    // What-if calculator
    monthlyIncome = 5000;
    savingsPercent = 20;
    emergencyMonths = 6;

    constructor() {
        this.goalForm = this.fb.group({
            name: ['', Validators.required],
            category: ['emergency', Validators.required],
            target_amount: [1000, [Validators.required, Validators.min(1)]],
            deadline: ['']
        });
    }

    ngOnInit(): void {
        this.loadGoals();
    }

    loadGoals(): void {
        this.budgetService.getGoals().subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.goals.set(res.data || []);
                }
                this.loading.set(false);
            },
            error: () => this.loading.set(false)
        });
    }

    createGoal(): void {
        if (this.goalForm.invalid) return;
        this.budgetService.createGoal(this.goalForm.value).subscribe({
            next: (res: any) => {
                if (res.success) {
                    this.goals.update(g => [res.data, ...g]);
                    this.showGoalForm.set(false);
                    this.goalForm.reset({ category: 'emergency', target_amount: 1000 });
                }
            }
        });
    }

    executeRoundup(): void {
        this.budgetService.executeRoundup().subscribe({
            next: (res: any) => {
                if (res.success) {
                    alert(`Roundup executed! Saved $${res.data.total_saved} from ${res.data.transactions_rounded} transactions.`);
                    this.loadGoals();
                }
            }
        });
    }

    getProgress(goal: any): number {
        const current = parseFloat(goal.current_amount) || 0;
        const target = parseFloat(goal.target_amount) || 1;
        return Math.min((current / target) * 100, 100);
    }

    getCategoryIcon(cat: string): string {
        return this.categories.find(c => c.value === cat)?.icon || '🎯';
    }

    get monthlySavings(): number { return this.monthlyIncome * (this.savingsPercent / 100); }
    get emergencyTarget(): number { return (this.monthlyIncome - this.monthlySavings) * this.emergencyMonths; }
    get monthsToEmergency(): number { return this.monthlySavings > 0 ? Math.ceil(this.emergencyTarget / this.monthlySavings) : 0; }

    Math = Math;
}
