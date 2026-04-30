import { Component, OnInit, ViewChild, ElementRef, AfterViewInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { createChart, LineSeries, Time } from 'lightweight-charts';
import { TradingService } from '../../core/services/trading.service';
import { StockService } from '../../core/services/stock.service';
import { firstValueFrom } from 'rxjs';

@Component({
  selector: 'app-trading',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './trading.component.html',
  styleUrls: ['./trading.component.scss']
})
export class TradingComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('chartContainer') chartContainer!: ElementRef;
  private chart!: any;
  private lineSeries!: any;
  Math = Math;


  isVerified = false;
  showKycModal = false;
  kycForm: FormGroup;

  marketSummary: any[] = [];
  portfolio: any = null;
  buyingPower: number = 0;
  totalPortfolioValue: number = 0;

  // Search & Active Stock
  searchForm: FormGroup;
  searchResults: any[] = [];
  activeStock: any = null;
  activeStockDetails: any = null;
  loadingStock = false;
  searchLoading = false;
  searchMessage = '';

  // Trade Panel
  tradeForm: FormGroup;
  tradeType: 'buy' | 'sell' = 'buy';
  orderProcessing = false;
  quickQuantities = [1, 5, 10, 25];
  selectedRange = '1M';
  timeRanges = ['1D', '1W', '1M', '1Y'];

  // Trade notification
  tradeNotification: { type: 'success' | 'error'; message: string } | null = null;

  watchlist: string[] = ['AAPL', 'TSLA', 'GOOGL', 'AMZN', 'MSFT'];
  watchlistData: any[] = [];

  // Order history
  orderHistory: any[] = [];
  showOrderHistory = false;

  // Category filter
  activeCategory = 'All';
  categories = ['All', 'Technology', 'Finance', 'Healthcare', 'Consumer', 'Crypto'];

  // Live simulation
  private simulationInterval: any = null;
  priceFlashState: 'up' | 'down' | null = null;
  flashTimeout: any = null;
  lastTickAt: Date | null = null;

  constructor(
    private fb: FormBuilder,
    private tradingService: TradingService,
    private stockService: StockService
  ) {
    this.kycForm = this.fb.group({
      dob: ['', Validators.required],
      ssn: ['', [Validators.required, Validators.pattern('^[0-9]{9}$')]],
      sec_q1: ['', Validators.required],
      sec_a1: ['', Validators.required],
      sec_q2: ['', Validators.required],
      sec_a2: ['', Validators.required],
      risk_ack: [false, Validators.requiredTrue]
    });

    this.searchForm = this.fb.group({
      query: ['']
    });

    this.tradeForm = this.fb.group({
      quantity: [1, [Validators.required, Validators.min(1)]]
    });
  }

  async ngOnInit() {
    await this.checkKycStatus();
    await this.loadDashboard();
    this.startLiveSimulation();
  }

  ngAfterViewInit() {
    this.initChart();
  }

  ngOnDestroy() {
    if (this.chart) {
      this.chart.remove();
    }
    if (this.simulationInterval) {
      clearInterval(this.simulationInterval);
    }
    if (this.flashTimeout) {
      clearTimeout(this.flashTimeout);
    }
  }

  async checkKycStatus() {
    try {
      const res = await firstValueFrom(this.tradingService.getAccount());
      if (res.data && res.data.status === 'verified') {
        this.isVerified = true;
      }
    } catch (e) {
      this.isVerified = false;
    }
  }

  async loadDashboard() {
    if (this.isVerified) {
       await this.loadPortfolio();
       await this.loadOrderHistory();
    }
    try {
      const marketRes = await firstValueFrom(this.stockService.getMarketSummary());
      if (marketRes.data) this.marketSummary = marketRes.data;
    } catch(e) {}

    await this.loadWatchlistData();
    // Load default stock
    await this.selectStock('AAPL');
  }

  async loadPortfolio() {
    if (!this.isVerified) return;
    try {
      const res = await firstValueFrom(this.tradingService.getPortfolio());
      if (res.data) {
        this.portfolio = res.data.positions || [];
        this.buyingPower = res.data.buying_power || 0;
        this.calculatePortfolioValue();
      }
    } catch(e) {
      console.error("Failed to load portfolio", e);
    }
  }

  async loadOrderHistory() {
    if (!this.isVerified) return;
    try {
      const res = await firstValueFrom(this.tradingService.getOrderHistory());
      this.orderHistory = res.data || [];
    } catch(e) {}
  }

  calculatePortfolioValue() {
    let posValue = 0;
    if (this.portfolio) {
      for(let p of this.portfolio) {
        const currentPrice = this.activeStock && this.activeStock.symbol === p.symbol ? this.activeStock.price : p.avg_cost;
        posValue += (p.quantity * currentPrice);
      }
    }
    this.totalPortfolioValue = posValue + this.buyingPower;
  }

  async loadWatchlistData() {
    this.watchlistData = [];
    for (const sym of this.watchlist) {
      try {
        const q = await firstValueFrom(this.stockService.getQuote(sym));
        if(q.data) {
          this.watchlistData.push({
            symbol: sym,
            price: q.data.price,
            change: q.data.change_percent,
            flash: null
          });
        }
      } catch(e) {}
    }
  }

  async submitKyc() {
    if (this.kycForm.invalid) return;
    try {
      await firstValueFrom(this.tradingService.verifyAccount(this.kycForm.value));
      this.isVerified = true;
      this.showKycModal = false;
      this.showTradeNotification('success', 'Trading account verified. Paper buying power is ready.');
      await this.loadPortfolio();
      await this.loadOrderHistory();
    } catch (e) {
      this.showTradeNotification('error', 'Verification failed. Check the form and try again.');
    }
  }

  async searchStock() {
    const query = this.searchForm.value.query;
    if (!query) return;
    this.searchLoading = true;
    this.searchMessage = '';
    try {
      const res = await firstValueFrom(this.stockService.search(query));
      this.searchResults = res.data || [];
      if (!this.searchResults.length) {
        this.searchMessage = 'No matches found';
      }
    } catch(e) {
      this.searchResults = [];
      this.searchMessage = 'Search is temporarily unavailable';
    } finally {
      this.searchLoading = false;
    }
  }

  filterByCategory(category: string) {
    this.activeCategory = category;
    if (category === 'All') {
      this.searchResults = [];
      this.searchForm.patchValue({query: ''});
    } else {
      this.searchForm.patchValue({query: category});
      this.searchStock();
    }
  }

  async selectStock(symbol: string) {
    this.searchResults = [];
    this.searchForm.patchValue({query: ''});
    this.loadingStock = true;
    this.activeStockDetails = null;

    try {
      const quoteRes = await firstValueFrom(this.stockService.getQuote(symbol));
      this.activeStock = quoteRes.data;

      if(!this.chart) {
        this.initChart();
      }

      const stockServiceAny = this.stockService as any;
      if (typeof stockServiceAny.getDetails === 'function') {
        stockServiceAny.getDetails(symbol).subscribe((res: any) => {
           this.activeStockDetails = res.data;
        });
      }

      const chartRes = await firstValueFrom(this.stockService.getChart(symbol));
      if (chartRes.data && this.lineSeries) {
        const d = chartRes.data as any[];
        const formatted = d.map(item => ({...item, time: this.formatTime(item.time)}));
        this.lineSeries.setData(formatted);
      }

      this.tradeForm.reset({quantity: 1});
      this.tradeType = 'buy';
    } catch(e) {
      console.warn("Failed to select stock", e);
      this.showTradeNotification('error', `Could not load ${symbol}. Try another symbol.`);
    } finally {
      this.loadingStock = false;
    }
  }

  formatTime(timeStr: string | number): Time {
    if (typeof timeStr === 'number') return timeStr as Time;
    return (new Date(timeStr).getTime() / 1000) as Time;
  }

  initChart() {
    if (!this.chartContainer) return;
    this.chart = createChart(this.chartContainer.nativeElement, {
      width: this.chartContainer.nativeElement.clientWidth,
      height: 480,
      layout: {
        background: { color: 'transparent' },
        textColor: '#94a3b8',
      },
      grid: {
        vertLines: { color: 'rgba(31, 41, 55, 0.4)', style: 1 },
        horzLines: { color: 'rgba(31, 41, 55, 0.4)', style: 1 },
      },
      timeScale: {
        borderColor: '#1f2937',
        timeVisible: true,
        secondsVisible: false,
      },
      rightPriceScale: {
        borderColor: '#1f2937',
      },
      crosshair: {
        mode: 0,
      }
    });

    this.lineSeries = this.chart.addSeries(LineSeries, {
      color: '#3b82f6',
      lineWidth: 2,
      crosshairMarkerVisible: true,
      crosshairMarkerRadius: 5,
      crosshairMarkerBorderColor: '#60a5fa',
      crosshairMarkerBackgroundColor: '#1d4ed8'
    });

    new ResizeObserver(entries => {
      if (entries.length === 0 || entries[0].target !== this.chartContainer.nativeElement) { return; }
      const newRect = entries[0].contentRect;
      this.chart.applyOptions({ height: newRect.height, width: newRect.width });
    }).observe(this.chartContainer.nativeElement);
  }

  setTradeType(type: 'buy' | 'sell') {
    this.tradeType = type;
  }

  setQuantity(quantity: number) {
    this.tradeForm.patchValue({ quantity });
  }

  setMaxQuantity() {
    const max = this.maxTradableQuantity();
    if (max > 0) {
      this.tradeForm.patchValue({ quantity: max });
    }
  }

  setRange(range: string) {
    this.selectedRange = range;
  }

  async executeTrade() {
    if (!this.canSubmitOrder()) return;

    if (!this.isVerified) {
       this.showKycModal = true;
       return;
    }

    this.orderProcessing = true;
    const payload = {
      symbol: this.activeStock.symbol,
      type: this.tradeType,
      quantity: this.tradeForm.value.quantity,
      price: this.activeStock.price
    };

    try {
      await firstValueFrom(this.tradingService.executeTrade(payload));
      this.showTradeNotification('success', `${this.tradeType.toUpperCase()} order executed: ${payload.quantity} ${payload.symbol} @ $${payload.price.toFixed(2)}`);
      this.tradeForm.reset({quantity: 1});
      await this.loadPortfolio();
      await this.loadOrderHistory();
    } catch (e: any) {
      this.showTradeNotification('error', e.error?.message || 'Trade failed. Check your balance.');
    } finally {
      this.orderProcessing = false;
    }
  }

  showTradeNotification(type: 'success' | 'error', message: string) {
    this.tradeNotification = { type, message };
    setTimeout(() => { this.tradeNotification = null; }, 5000);
  }

  getPositionPL(position: any): number {
    if (!this.activeStock || this.activeStock.symbol !== position.symbol) return 0;
    return (this.activeStock.price - position.avg_cost) * position.quantity;
  }

  getPositionMarketPrice(position: any): number {
    const watch = this.watchlistData.find(w => w.symbol === position.symbol);
    if (this.activeStock?.symbol === position.symbol) return this.activeStock.price;
    return watch?.price || position.avg_cost;
  }

  getPositionMarketValue(position: any): number {
    return this.getPositionMarketPrice(position) * position.quantity;
  }

  getPositionPLPercent(position: any): number {
    const cost = position.avg_cost * position.quantity;
    if (!cost) return 0;
    return (this.getPositionPL(position) / cost) * 100;
  }

  addToWatchlist(symbol: string) {
    if (!this.watchlist.includes(symbol)) {
      this.watchlist.push(symbol);
      this.loadWatchlistData();
      this.showTradeNotification('success', `${symbol} added to watchlist`);
    } else {
      this.showTradeNotification('error', `${symbol} is already on your watchlist`);
    }
  }

  removeFromWatchlist(symbol: string, event?: Event) {
    event?.stopPropagation();
    this.watchlist = this.watchlist.filter(s => s !== symbol);
    this.watchlistData = this.watchlistData.filter(w => w.symbol !== symbol);
  }

  estimatedCost(): number {
    return (Number(this.tradeForm.value.quantity) || 0) * (this.activeStock?.price || 0);
  }

  maxTradableQuantity(): number {
    if (!this.activeStock?.price) return 0;
    if (this.tradeType === 'buy') {
      return Math.floor(this.buyingPower / this.activeStock.price);
    }
    const position = this.portfolio?.find((p: any) => p.symbol === this.activeStock.symbol);
    return position?.quantity || 0;
  }

  canSubmitOrder(): boolean {
    return !!this.activeStock && !this.orderProcessing && this.tradeForm.valid && !this.orderWarning();
  }

  orderWarning(): string {
    if (!this.activeStock || this.tradeForm.invalid) return '';
    const quantity = Number(this.tradeForm.value.quantity) || 0;
    if (quantity <= 0) return 'Enter a quantity greater than zero.';
    if (this.tradeType === 'buy' && this.estimatedCost() > this.buyingPower) {
      return 'Estimated cost is above your buying power.';
    }
    if (this.tradeType === 'sell' && quantity > this.maxTradableQuantity()) {
      return 'You do not own enough shares to place this sell order.';
    }
    return '';
  }

  riskLevel(): { label: string; className: string } {
    const change = Math.abs(this.activeStock?.change_percent || 0);
    if (change >= 3) return { label: 'High volatility', className: 'high' };
    if (change >= 1) return { label: 'Moderate volatility', className: 'medium' };
    return { label: 'Low volatility', className: 'low' };
  }

  portfolioInvestedValue(): number {
    return (this.portfolio || []).reduce((sum: number, p: any) => sum + this.getPositionMarketValue(p), 0);
  }

  dayChangeEstimate(): number {
    const invested = this.portfolioInvestedValue();
    const blendedChange = this.watchlistData.length
      ? this.watchlistData.reduce((sum, w) => sum + (w.change || 0), 0) / this.watchlistData.length
      : 0.24;
    return invested * (blendedChange / 100);
  }

  diversificationScore(): number {
    const symbols = new Set((this.portfolio || []).map((p: any) => p.symbol));
    return Math.min(100, 35 + symbols.size * 13 + this.allocation().length * 6);
  }

  buyingPowerUsage(): number {
    const total = this.totalPortfolioValue || 1;
    return Math.min(100, (this.portfolioInvestedValue() / total) * 100);
  }

  allocation(): any[] {
    const colors = ['#3b82f6', '#22c55e', '#f59e0b', '#ec4899', '#8b5cf6', '#06b6d4'];
    const sectorTotals = new Map<string, number>();
    for (const p of this.portfolio || []) {
      const sector = this.sectorForSymbol(p.symbol);
      sectorTotals.set(sector, (sectorTotals.get(sector) || 0) + this.getPositionMarketValue(p));
    }
    const total = Array.from(sectorTotals.values()).reduce((sum, value) => sum + value, 0);
    if (!total) {
      return [
        { sector: 'Cash', value: this.buyingPower || 0, percent: 100, color: colors[0] }
      ];
    }
    return Array.from(sectorTotals.entries()).map(([sector, value], index) => ({
      sector,
      value,
      percent: Math.round((value / total) * 100),
      color: colors[index % colors.length]
    }));
  }

  pieSegments(): any[] {
    let offset = 0;
    return this.allocation().map(slice => {
      const dash = Math.max(0, (slice.percent / 100) * 503);
      const segment = { ...slice, dash, offset: -offset };
      offset += dash;
      return segment;
    });
  }

  rebalanceSuggestions(): any[] {
    return this.allocation()
      .filter(slice => Math.abs(slice.percent - this.targetAllocation(slice.sector)) >= 4)
      .map(slice => {
        const target = this.targetAllocation(slice.sector);
        return {
          sector: slice.sector,
          current: slice.percent,
          target,
          action: slice.percent > target ? 'Trim' : 'Add',
          delta: Math.abs(slice.percent - target)
        };
      })
      .slice(0, 3);
  }

  targetAllocation(sector: string): number {
    const targets: Record<string, number> = {
      Technology: 30,
      Finance: 18,
      Healthcare: 18,
      Consumer: 16,
      Crypto: 8,
      Cash: 10
    };
    return targets[sector] || 10;
  }

  sentimentItems(): any[] {
    const symbol = this.activeStock?.symbol || 'AAPL';
    const sector = this.activeStockDetails?.sector || this.sectorForSymbol(symbol);
    return [
      { tone: 'bullish', label: 'Bullish', headline: `${symbol} momentum improves as volume trends above baseline`, source: 'FlowPay Signals' },
      { tone: 'neutral', label: 'Neutral', headline: `${sector} peers remain mixed while rates stay in focus`, source: 'Market Desk' },
      { tone: this.riskLevel().className === 'high' ? 'bearish' : 'bullish', label: this.riskLevel().className === 'high' ? 'Bearish' : 'Bullish', headline: `${this.riskLevel().label} detected on today's tape`, source: 'Risk Engine' }
    ];
  }

  taxSummary() {
    const gains = Math.max(0, this.portfolioInvestedValue() * 0.045);
    const losses = Math.max(0, this.portfolioInvestedValue() * 0.012);
    const net = gains - losses;
    return {
      gains,
      losses,
      net,
      estimatedTax: Math.max(0, net * 0.15)
    };
  }

  marketStatus(): string {
    const now = new Date();
    const day = now.getDay();
    const minutes = now.getHours() * 60 + now.getMinutes();
    const open = 9 * 60 + 30;
    const close = 16 * 60;
    if (day === 0 || day === 6) return 'Market closed';
    if (minutes >= open && minutes < close) return 'Market open';
    return 'After hours';
  }

  sectorForSymbol(symbol: string): string {
    const sectors: Record<string, string> = {
      AAPL: 'Technology', MSFT: 'Technology', GOOGL: 'Technology', AMZN: 'Technology', NVDA: 'Technology', META: 'Technology',
      JPM: 'Finance', V: 'Finance', GS: 'Finance',
      JNJ: 'Healthcare', UNH: 'Healthcare', PFE: 'Healthcare',
      TSLA: 'Consumer', DIS: 'Consumer', NKE: 'Consumer',
      BTC: 'Crypto', ETH: 'Crypto', SOL: 'Crypto'
    };
    return sectors[symbol] || 'Other';
  }

  // --- Live Data Simulator ---
  startLiveSimulation() {
    this.simulationInterval = setInterval(() => {
      this.simulateMarketTick();
    }, 2000);
  }

  simulateMarketTick() {
    if (this.activeStock) {
      const volatility = 0.001;
      const move = this.activeStock.price * volatility * (Math.random() - 0.48);
      const newPrice = this.activeStock.price + move;

      if(this.flashTimeout) clearTimeout(this.flashTimeout);
      this.priceFlashState = move >= 0 ? 'up' : 'down';
      this.flashTimeout = setTimeout(() => { this.priceFlashState = null; }, 1000);

      this.activeStock.price = newPrice;
      this.activeStock.change = (this.activeStock.change || 0) + move;
      this.lastTickAt = new Date();

      if (this.lineSeries) {
         try {
           this.lineSeries.update({
             time: (Date.now() / 1000) as Time,
             value: newPrice
           });
         } catch(e) {}
      }
    }

    if (this.marketSummary) {
      this.marketSummary.forEach(m => {
        if(Math.random() > 0.5) {
          const move = (Math.random() - 0.5) * 10;
          m.price += move;
          m.change_percent += (move / 100);
        }
      });
    }

    if (this.watchlistData) {
      this.watchlistData.forEach(w => {
         if(Math.random() > 0.3) {
            const move = w.price * 0.001 * (Math.random() - 0.5);
            w.price += move;
            w.change += (move / w.price) * 100;
            w.flash = move >= 0 ? 'up' : 'down';
            setTimeout(() => { w.flash = null; }, 1000);
         }
      });
    }

    this.calculatePortfolioValue();
  }
}
