import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class SubscriptionService {
    private apiUrl = `${environment.apiUrl}/subscriptions`;

    constructor(private http: HttpClient) { }

    getSubscriptions() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}`);
    }

    trackSubscription(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/track`, data);
    }

    setAutoPay(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/autopay`, data);
    }
}
