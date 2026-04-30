import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class BudgetService {
    private budgetUrl = `${environment.apiUrl}/budget`;
    private autosaveUrl = `${environment.apiUrl}/autosave`;

    constructor(private http: HttpClient) { }

    getGoals() {
        return this.http.get<ApiResponse<any>>(`${this.budgetUrl}/goals`);
    }

    createGoal(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.budgetUrl}/goals`, data);
    }

    createAutoSaveRule(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.autosaveUrl}/rules`, data);
    }

    executeRoundup() {
        return this.http.post<ApiResponse<any>>(`${this.autosaveUrl}/roundup/execute`, {});
    }
}
