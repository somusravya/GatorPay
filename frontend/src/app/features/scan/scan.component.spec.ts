import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ScanComponent } from './scan.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { QrService } from '../../core/services/qr.service';
import { of, throwError } from 'rxjs';

describe('ScanComponent', () => {
  let component: ScanComponent;
  let fixture: ComponentFixture<ScanComponent>;

  beforeEach(async () => {
    const mockQrService = {
      registerMerchant: jasmine.createSpy().and.returnValue(of({ success: true })),
      generateQR: jasmine.createSpy().and.returnValue(of({ success: true, data: { base64_png: 'data:img' } })),
      payQR: jasmine.createSpy().and.returnValue(of({ success: true }))
    };

    await TestBed.configureTestingModule({
      imports: [ScanComponent, HttpClientTestingModule],
      providers: [
        { provide: QrService, useValue: mockQrService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ScanComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
