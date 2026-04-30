import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';

export interface SavingsGoal {
    id: string;
    name: string;
    targetAmount: number;
    savedAmount: number;
    icon: string;
    color: string;
    deadline: string;
    createdAt: string;
}

@Injectable({ providedIn: 'root' })
export class SavingsGoalsService {
    private readonly storageKey = 'gatorpay_savings_goals';

    getGoals(): Observable<SavingsGoal[]> {
        return of(this.readGoals());
    }

    createGoal(goal: Omit<SavingsGoal, 'id' | 'savedAmount' | 'createdAt'>): Observable<SavingsGoal> {
        const goals = this.readGoals();
        const created: SavingsGoal = {
            ...goal,
            id: crypto.randomUUID?.() ?? Date.now().toString(),
            savedAmount: 0,
            createdAt: new Date().toISOString()
        };
        goals.push(created);
        this.saveGoals(goals);
        return of(created);
    }

    addMoney(goalId: string, amount: number): Observable<SavingsGoal | null> {
        const goals = this.readGoals();
        const goal = goals.find(g => g.id === goalId);
        if (!goal) return of(null);

        goal.savedAmount = Math.min(goal.savedAmount + amount, goal.targetAmount);
        this.saveGoals(goals);
        return of(goal);
    }

    deleteGoal(goalId: string): Observable<boolean> {
        const goals = this.readGoals().filter(g => g.id !== goalId);
        this.saveGoals(goals);
        return of(true);
    }

    private readGoals(): SavingsGoal[] {
        const raw = localStorage.getItem(this.storageKey);
        if (!raw) return [];
        try {
            return JSON.parse(raw) as SavingsGoal[];
        } catch {
            return [];
        }
    }

    private saveGoals(goals: SavingsGoal[]): void {
        localStorage.setItem(this.storageKey, JSON.stringify(goals));
    }
}
