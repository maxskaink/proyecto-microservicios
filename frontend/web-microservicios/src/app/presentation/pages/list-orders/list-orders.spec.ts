import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListOrders } from './list-orders';

describe('ListOrders', () => {
  let component: ListOrders;
  let fixture: ComponentFixture<ListOrders>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListOrders]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ListOrders);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
