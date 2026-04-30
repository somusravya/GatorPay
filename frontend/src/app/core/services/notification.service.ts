import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class NotificationService {
    private apiUrl = `${environment.apiUrl}/notifications`;

    constructor(private http: HttpClient) { }

    getNotifications(filterType?: string) {
        const params = filterType ? `?type=${filterType}` : '';
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}${params}`);
    }

    markRead(id: string) {
        return this.http.put<ApiResponse<any>>(`${this.apiUrl}/${id}/read`, {});
    }

    markAsRead(id: string) {
        return this.markRead(id);
    }

    markAllRead() {
        return this.http.put<ApiResponse<any>>(`${this.apiUrl}/all/read`, {});
    }

    getPreferences() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/preferences`);
    }

    updatePreferences(prefs: any) {
        return this.http.put<ApiResponse<any>>(`${this.apiUrl}/preferences`, prefs);
    }
}
