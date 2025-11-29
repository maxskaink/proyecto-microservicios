import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShippingOrderPage } from './shipping-order-page';

describe('ShippingOrderPage', () => {
  let component: ShippingOrderPage;
  let fixture: ComponentFixture<ShippingOrderPage>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ShippingOrderPage]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShippingOrderPage);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
