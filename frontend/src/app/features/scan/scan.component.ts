import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { QrService } from '../../core/services/qr.service';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-scan',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './scan.component.html',
  styleUrls: ['./scan.component.scss']
})
export class ScanComponent {
  activeTab: 'pay' | 'receive' = 'pay';
  
  // Pay Tab
  scanCode: string = '';
  payAmount: number = 0;
  
  // Receive Tab
  isMerchant: boolean = false;
  businessName: string = '';
  category: string = '';
  qrCodeImage: string = '';
  qrAmount: number = 0;
  isDynamic: boolean = false;
  
  constructor(private qrService: QrService) {}

  async payQR() {
    if (!this.scanCode) return;
    try {
      await firstValueFrom(this.qrService.payQR(this.scanCode, this.payAmount));
      alert("Payment successful! You earned 1.5% cashback.");
      this.scanCode = '';
      this.payAmount = 0;
    } catch (e: any) {
      alert("Payment failed: " + e.error?.message);
    }
  }

  async registerMerchant() {
    try {
      await firstValueFrom(this.qrService.registerMerchant({
        business_name: this.businessName,
        category: this.category
      }));
      this.isMerchant = true;
      this.generateQR();
    } catch (e: any) {
      alert("Merchant registration failed.");
    }
  }

  async generateQR() {
    try {
      const res = await firstValueFrom(this.qrService.generateQR(this.qrAmount, this.isDynamic));
      if (res.data?.base64_png) {
        this.qrCodeImage = res.data.base64_png;
        this.isMerchant = true; // Auto mark as merchant if it succeeded
      }
    } catch (e: any) {
      if (e.status === 400 || e.status === 404) {
        // Not a merchant yet
        this.isMerchant = false;
      }
    }
  }
}
