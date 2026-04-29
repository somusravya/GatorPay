import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class AdminService {
    private apiUrl = `${environment.apiUrl}/admin`;

    constructor(private http: HttpClient) { }

    getMetrics() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/metrics`);
    }

    getUsers(page: number = 1, limit: number = 20, search?: string) {
        const params = search ? `&search=${search}` : '';
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/users?page=${page}&limit=${limit}${params}`);
    }

    getFraudReview() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/fraud/review`);
    }
}
