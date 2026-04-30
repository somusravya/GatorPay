import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { QrService } from '../../core/services/qr.service';
import { firstValueFrom } from 'rxjs';

declare var jsQR: any;

@Component({
  selector: 'app-scan',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './scan.component.html',
  styleUrls: ['./scan.component.scss']
})
export class ScanComponent {
  activeTab: 'pay' | 'receive' = 'pay';

  // Payment Flow Steps
  payStep: 'scan' | 'verify' | 'amount' | 'confirm' | 'result' = 'scan';
  scanCode: string = '';
  payAmount: number = 0;
  recipientInfo: any = null;
  paymentResult: 'success' | 'error' | null = null;
  paymentError: string = '';
  isProcessing = false;

  // QR Decode
  isDragging = false;
  decodeError = '';

  // Receive Tab
  isMerchant: boolean = false;
  businessName: string = '';
  category: string = '';
  qrCodeImage: string = '';
  qrAmount: number = 0;
  isDynamic: boolean = false;

  // Recent payments (local tracking)
  recentPayments: { to: string; amount: number; time: Date; status: string }[] = [];

  // Toast
  toast: { message: string; type: 'success' | 'error' | 'info' } | null = null;

  constructor(private qrService: QrService) {}

  // --- QR File Upload and Decode ---
  onFileSelected(event: any) {
    const file = event.target.files?.[0];
    if (file) {
      this.decodeQRFromFile(file);
    }
  }

  onDragOver(event: DragEvent) {
    event.preventDefault();
    this.isDragging = true;
  }

  onDragLeave(event: DragEvent) {
    event.preventDefault();
    this.isDragging = false;
  }

  onDrop(event: DragEvent) {
    event.preventDefault();
    this.isDragging = false;
    const file = event.dataTransfer?.files[0];
    if (file) {
      this.decodeQRFromFile(file);
    }
  }

  async decodeQRFromFile(file: File) {
    this.decodeError = '';
    const img = new Image();
    const url = URL.createObjectURL(file);

    img.onload = () => {
      const canvas = document.createElement('canvas');
      canvas.width = img.width;
      canvas.height = img.height;
      const ctx = canvas.getContext('2d');
      if (!ctx) {
        this.decodeError = 'Canvas not supported in this browser.';
        return;
      }
      ctx.drawImage(img, 0, 0);
      const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);

      try {
        // Dynamic import for jsQR
        const jsQRModule = (window as any).jsQR || null;
        if (jsQRModule) {
          const code = jsQRModule(imageData.data, imageData.width, imageData.height);
          if (code && code.data) {
            this.scanCode = code.data;
            this.showToast('QR code decoded successfully!', 'success');
            this.lookupQR();
          } else {
            this.decodeError = 'Could not read QR code from this image. Try a clearer image.';
          }
        } else {
          // Fallback: just prompt for manual entry
          this.decodeError = 'QR decoder not available. Please enter the code manually.';
        }
      } catch (e) {
        this.decodeError = 'Error decoding QR code. Please enter the code manually.';
      }
      URL.revokeObjectURL(url);
    };

    img.onerror = () => {
      this.decodeError = 'Could not load image file.';
      URL.revokeObjectURL(url);
    };

    img.src = url;
  }

  // --- Payment Flow ---
  async lookupQR() {
    if (!this.scanCode) return;
    this.decodeError = '';
    this.isProcessing = true;

    try {
      const res = await firstValueFrom(this.qrService.lookupQR(this.scanCode));
      this.recipientInfo = res.data;
      this.payAmount = parseFloat(this.recipientInfo?.amount) || 0;
      this.payStep = 'verify';
    } catch (e: any) {
      this.decodeError = e.error?.message || 'Invalid or unrecognized QR code.';
      this.recipientInfo = null;
    } finally {
      this.isProcessing = false;
    }
  }

  proceedToAmount() {
    this.payStep = 'amount';
  }

  proceedToConfirm() {
    if (this.payAmount <= 0) {
      this.showToast('Please enter a valid amount.', 'error');
      return;
    }
    this.payStep = 'confirm';
  }

  async confirmPayment() {
    this.isProcessing = true;
    this.paymentResult = null;
    this.paymentError = '';

    try {
      await firstValueFrom(this.qrService.payQR(this.scanCode, this.payAmount));
      this.paymentResult = 'success';
      this.payStep = 'result';
      this.recentPayments.unshift({
        to: this.recipientInfo?.business_name || 'Unknown',
        amount: this.payAmount,
        time: new Date(),
        status: 'success'
      });
    } catch (e: any) {
      this.paymentResult = 'error';
      this.paymentError = e.error?.message || 'Payment failed. Please try again.';
      this.payStep = 'result';
      this.recentPayments.unshift({
        to: this.recipientInfo?.business_name || 'Unknown',
        amount: this.payAmount,
        time: new Date(),
        status: 'failed'
      });
    } finally {
      this.isProcessing = false;
    }
  }

  resetPayment() {
    this.payStep = 'scan';
    this.scanCode = '';
    this.payAmount = 0;
    this.recipientInfo = null;
    this.paymentResult = null;
    this.paymentError = '';
    this.decodeError = '';
  }

  // --- Receive Tab ---
  async registerMerchant() {
    try {
      await firstValueFrom(this.qrService.registerMerchant({
        business_name: this.businessName,
        category: this.category
      }));
      this.isMerchant = true;
      this.showToast('Merchant registered!', 'success');
      this.generateQR();
    } catch (e: any) {
      this.showToast('Merchant registration failed.', 'error');
    }
  }

  async generateQR() {
    try {
      const res = await firstValueFrom(this.qrService.generateQR(this.qrAmount, this.isDynamic));
      if (res.data?.base64_png) {
        this.qrCodeImage = res.data.base64_png;
        this.isMerchant = true;
      }
    } catch (e: any) {
      if (e.status === 400 || e.status === 404) {
        this.isMerchant = false;
      }
    }
  }

  showToast(message: string, type: 'success' | 'error' | 'info') {
    this.toast = { message, type };
    setTimeout(() => { this.toast = null; }, 4000);
  }
}
