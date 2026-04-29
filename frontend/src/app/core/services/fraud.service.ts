import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class FraudService {
    private apiUrl = `${environment.apiUrl}/fraud`;

    constructor(private http: HttpClient) { }

    getAlerts(status?: string) {
        const params = status ? `?status=${status}` : '';
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/alerts${params}`);
    }

    reviewAlert(data: { alert_id: string; action: string }) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/review`, data);
    }
}
