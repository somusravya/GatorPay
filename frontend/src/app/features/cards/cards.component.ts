import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { CardService } from '../../core/services/card.service';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-cards',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './cards.component.html',
  styleUrls: ['./cards.component.scss']
})
export class CardsComponent implements OnInit {
  cards: any[] = [];
  
  // Create Card
  showCreateModal = false;
  newCardName = '';
  
  // Show Details (OTP Flow)
  showOtpModal = false;
  targetCardId: string | null = null;
  otpInput = '';
  unlockedDetails: {[id: string]: any} = {};

  constructor(private cardService: CardService) {}

  ngOnInit() {
    this.loadCards();
  }

  async loadCards() {
    try {
      const res = await firstValueFrom(this.cardService.getCards());
      this.cards = res.data || [];
    } catch(e) {}
  }

  async createCard() {
    if (!this.newCardName) return;
    try {
      await firstValueFrom(this.cardService.createCard(this.newCardName));
      this.newCardName = '';
      this.showCreateModal = false;
      this.loadCards();
    } catch (e: any) {
      alert("Failed to create card: " + e.error?.message);
    }
  }

  async toggleFreeze(cardId: string) {
    try {
      await firstValueFrom(this.cardService.freezeCard(cardId));
      this.loadCards();
    } catch (e) {
      alert("Failed to update status");
    }
  }

  async requestOTP(cardId: string) {
    this.targetCardId = cardId;
    try {
      await firstValueFrom(this.cardService.requestOTP(cardId));
      this.showOtpModal = true;
    } catch(e) {
      alert("Failed to send OTP");
    }
  }

  async verifyAndShowDetails() {
    if (!this.targetCardId || !this.otpInput) return;
    try {
      const res = await firstValueFrom(this.cardService.getCardDetails(this.targetCardId, this.otpInput));
      this.unlockedDetails[this.targetCardId] = res.data;
      this.showOtpModal = false;
      this.otpInput = '';
    } catch (e: any) {
      alert("Invalid OTP");
    }
  }

  formatCardNumber(num: string): string {
    if (!num) return '';
    return num.match(/.{1,4}/g)?.join(' ') || num;
  }
}
