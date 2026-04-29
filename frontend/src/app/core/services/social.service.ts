import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { ApiResponse } from '../../shared/models';

@Injectable({ providedIn: 'root' })
export class SocialService {
    private apiUrl = `${environment.apiUrl}/social`;

    constructor(private http: HttpClient) { }

    getFeed() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/feed`);
    }

    createPost(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/post`, data);
    }

    reactToPost(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/react`, data);
    }

    getFriends() {
        return this.http.get<ApiResponse<any>>(`${this.apiUrl}/friends`);
    }

    addFriend(data: any) {
        return this.http.post<ApiResponse<any>>(`${this.apiUrl}/friends/add`, data);
    }
}
