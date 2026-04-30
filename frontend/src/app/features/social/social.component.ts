import { Component, OnInit, signal, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SocialService } from '../../core/services/social.service';

@Component({
    selector: 'app-social',
    standalone: true,
    imports: [CommonModule, FormsModule],
    templateUrl: './social.component.html',
    styleUrl: './social.component.scss'
})
export class SocialComponent implements OnInit {
    private socialService = inject(SocialService);
    feed = signal<any[]>([]);
    friends = signal<any[]>([]);
    loading = signal(true);
    newPost = { note: '', emoji: '💸', privacy: 'friends', amount: 0, recipient_id: '' };
    showComposer = false;
    showFriendDrawer = false;
    commentText: { [key: string]: string } = {};
    expandedComments: { [key: string]: boolean } = {};
    friendSearch = '';

    emojis = ['💸', '🎉', '🍕', '☕', '🎮', '🎵', '💪', '✈️', '🎁', '❤️'];

    ngOnInit(): void {
        this.socialService.getFeed().subscribe({
            next: (res: any) => { if (res.success) this.feed.set(res.data || []); this.loading.set(false); },
            error: () => this.loading.set(false)
        });
        this.socialService.getFriends().subscribe({
            next: (res: any) => { if (res.success) this.friends.set(res.data || []); }
        });
    }

    createPost(): void {
        if (!this.newPost.note.trim()) return;
        this.socialService.createPost(this.newPost).subscribe({
            next: (res: any) => {
                if (res.success) { this.feed.update(f => [res.data, ...f]); this.newPost = { note: '', emoji: '💸', privacy: 'friends', amount: 0, recipient_id: '' }; this.showComposer = false; }
            }
        });
    }

    likePost(item: any): void {
        this.socialService.reactToPost({ feed_item_id: item.id, type: 'like' }).subscribe({
            next: () => { item.like_count = (item.like_count || 0) + 1; }
        });
    }

    postComment(item: any): void {
        const content = this.commentText[item.id]?.trim();
        if (!content) return;
        this.socialService.reactToPost({ feed_item_id: item.id, type: 'comment', content }).subscribe({
            next: () => {
                item.comment_count = (item.comment_count || 0) + 1;
                this.commentText[item.id] = '';
            }
        });
    }

    toggleComments(itemId: string): void {
        this.expandedComments[itemId] = !this.expandedComments[itemId];
    }

    addFriend(friendId: string): void {
        this.socialService.addFriend({ friend_id: friendId }).subscribe({
            next: (res: any) => {
                if (res.success) { this.friends.update(f => [...f, res.data]); }
            }
        });
    }

    getTimeAgo(dateStr: string): string {
        const diff = Date.now() - new Date(dateStr).getTime();
        const mins = Math.floor(diff / 60000);
        if (mins < 60) return `${mins}m ago`;
        const hours = Math.floor(mins / 60);
        if (hours < 24) return `${hours}h ago`;
        return `${Math.floor(hours / 24)}d ago`;
    }
}
